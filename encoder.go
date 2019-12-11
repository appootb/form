package form

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
)

type Marshaler interface {
	MarshalURL() (string, error)
}

var (
	marshalerType = reflect.TypeOf((*Marshaler)(nil)).Elem()
)

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(src interface{}, dst url.Values) error {
	v := reflect.ValueOf(src)
	err := e.encode(v, dst)
	return err
}

// isValidStructPointer test if input value is a valid struct pointer.
func (e *Encoder) isValidStructPointer(v reflect.Value) bool {
	return v.Type().Kind() == reflect.Ptr && v.Elem().IsValid() && v.Elem().Type().Kind() == reflect.Struct
}

func (e *Encoder) isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func:
	case reflect.Map, reflect.Slice:
		return v.IsNil() || v.Len() == 0
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && e.isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && e.isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}

func (e *Encoder) encode(v reflect.Value, dst url.Values) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return errors.New("the interface must be a struct")
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		name, opts := fieldAlias(t.Field(i))
		if name == "-" {
			continue
		}

		fv := v.Field(i)

		// Encode struct pointer types if the field is a valid pointer and a struct.
		if e.isValidStructPointer(fv) {
			_ = e.encode(fv.Elem(), dst)
			continue
		}

		marshaler := e.getMarshaler(fv.Type(), fv)

		// Encode non-slice types and custom implementations immediately.
		if marshaler != nil {
			if opts.Contains("omitempty") && e.isZero(fv) {
				continue
			}

			value, err := marshaler.MarshalURL()
			if err != nil {
				return err
			}

			dst[name] = append(dst[name], value)
			continue
		}

		switch fv.Type().Kind() {
		case reflect.Struct:
			err := e.encode(fv, dst)
			if err != nil {
				return err
			}
		case reflect.Slice, reflect.Array:
			if fv.Len() == 0 && opts.Contains("omitempty") {
				continue
			}

			dst[name] = []string{}
			for j := 0; j < fv.Len(); j++ {
				value, err := e.getMarshaler(fv.Type().Elem(), fv.Index(j)).MarshalURL()
				if err != nil {
					return err
				}
				dst[name] = append(dst[name], value)
			}
		case reflect.Map:
			if fv.Len() == 0 && opts.Contains("omitempty") {
				continue
			}

			for _, k := range fv.MapKeys() {
				key, err := e.getMarshaler(k.Type(), k).MarshalURL()
				if err != nil {
					return err
				}
				value, err := e.getMarshaler(fv.MapIndex(k).Type(), fv.MapIndex(k)).MarshalURL()
				if err != nil {
					return err
				}
				dst[key] = append(dst[key], value)
			}
		default:
			return fmt.Errorf("marshaler not found for %v", fv.Type())
		}
	}

	return nil
}

func (e *Encoder) getMarshaler(t reflect.Type, v reflect.Value) Marshaler {
	if t.Implements(marshalerType) {
		return v.Interface().(Marshaler)
	}

	switch t.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return &Null{}
		}
		return e.getMarshaler(t.Elem(), v.Elem())
	case reflect.Bool:
		val := reflect.New(BoolType).Elem()
		val.SetBool(v.Bool())
		return val.Interface().(Marshaler)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := reflect.New(Int64Type).Elem()
		val.SetInt(v.Int())
		return val.Interface().(Marshaler)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := reflect.New(UInt64Type).Elem()
		val.SetUint(v.Uint())
		return val.Interface().(Marshaler)
	case reflect.Float32:
		val := reflect.New(Float32Type).Elem()
		val.SetFloat(v.Float())
		return val.Interface().(Marshaler)
	case reflect.Float64:
		val := reflect.New(Float64Type).Elem()
		val.SetFloat(v.Float())
		return val.Interface().(Marshaler)
	case reflect.String:
		val := reflect.New(StringType).Elem()
		val.SetString(v.String())
		return val.Interface().(Marshaler)
	case reflect.Interface:
		return &Interface{
			Val: v.Interface(),
		}
	//case reflect.Slice, reflect.Array:
	//case reflect.Map:
	//case reflect.Struct:
	//case reflect.Complex64, reflect.Complex128:
	default:
		return nil
	}
}
