package form

import (
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
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return TypeError
	}

	err := e.encode(v.Elem(), dst)
	return err
}

func (e *Encoder) isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func:
	case reflect.Map, reflect.Slice, reflect.Array:
		return v.IsNil() || v.Len() == 0
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
	var (
		marshaler Marshaler
	)

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		name, opts := fieldAlias(t.Field(i))
		if name == "-" {
			continue
		}

		fv := v.Field(i)
		if opts.Contains("omitempty") && e.isZero(fv) {
			continue
		}

		// Encode base types and custom implementations immediately.
		marshaler = e.getMarshaler(fv.Type(), fv)
		if marshaler != nil {
			value, err := marshaler.MarshalURL()
			if err != nil {
				return err
			}
			dst[name] = append(dst[name], value)
			continue
		}

		switch fv.Type().Kind() {
		case reflect.Ptr:
			if !fv.IsValid() || fv.IsNil() {
				dst[name] = []string{NullValue}
				continue
			}
			if err := e.encode(fv.Elem(), dst); err != nil {
				return err
			}
		case reflect.Struct:
			err := e.encode(fv, dst)
			if err != nil {
				return err
			}
		case reflect.Slice, reflect.Array:
			dst[name] = []string{}
			for j := 0; j < fv.Len(); j++ {
				value, err := e.getMarshaler(fv.Type().Elem(), fv.Index(j)).MarshalURL()
				if err != nil {
					return err
				}
				dst[name] = append(dst[name], value)
			}
		case reflect.Map:
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
	if v.CanAddr() && v.Addr().Type().Implements(marshalerType) {
		return v.Addr().Interface().(Marshaler)
	} else if v.Type().Implements(marshalerType) {
		if v.Type().Kind() != reflect.Ptr || v.IsValid() && !v.IsNil() {
			return v.Interface().(Marshaler)
		}
	}

	switch t.Kind() {
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
