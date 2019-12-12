package form

import (
	"fmt"
	"net/url"
	"reflect"
)

type Unmarshaler interface {
	UnmarshalURL(string) error
}

var (
	unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(dst interface{}, src url.Values) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return TypeError
	}

	fields := map[string]bool{}
	mapField, err := d.decode(v.Elem(), src, fields)
	if err != nil {
		return err
	}
	if !mapField.IsValid() {
		return nil
	}

	t := mapField.Type()
	m := reflect.MakeMapWithSize(reflect.MapOf(t.Key(), t.Elem()), len(src))
	for k, vals := range src {
		if fields[k] {
			continue
		}

		var (
			key = reflect.New(t.Key()).Elem()
			val = reflect.New(t.Elem()).Elem()
		)
		err := d.decodeElement(t.Key(), key, k)
		if err != nil {
			continue
		}
		err = d.decodeElement(t.Elem(), val, vals[0])
		if err != nil {
			val = reflect.Zero(t.Elem())
		}
		m.SetMapIndex(key, val)
	}

	mapField.Set(m)
	return nil
}

func (d *Decoder) decodeElement(t reflect.Type, v reflect.Value, src string) (err error) {
	if v.CanAddr() && v.Addr().Type().Implements(unmarshalerType) {
		err = v.Addr().Interface().(Unmarshaler).UnmarshalURL(src)
	} else if t.Implements(unmarshalerType) {
		err = v.Interface().(Unmarshaler).UnmarshalURL(src)
	} else {
		err = d.unmarshal(t, v, src)
	}
	return err
}

func (d *Decoder) decode(v reflect.Value, src url.Values, fields map[string]bool) (reflect.Value, error) {
	var (
		err            error
		mapField       reflect.Value
		recursionField reflect.Value
	)

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		name, _ := fieldAlias(t.Field(i))
		if name == "-" {
			fields[name] = false
			continue
		}

		fields[name] = true
		fv := v.Field(i)

		if fv.CanAddr() && fv.Addr().Type().Implements(unmarshalerType) {
			if err = fv.Addr().Interface().(Unmarshaler).UnmarshalURL(src.Get(name)); err != nil {
				goto End
			}
			continue
		}

		switch fv.Type().Kind() {
		case reflect.Ptr:
			if src.Get(name) == NullValue {
				fv.Set(reflect.Zero(fv.Type()))
				continue
			}
			fv.Set(reflect.New(fv.Type().Elem()))
			if fv.Type().Implements(unmarshalerType) {
				err = fv.Interface().(Unmarshaler).UnmarshalURL(src.Get(name))
			} else {
				recursionField, err = d.decode(fv.Elem(), src, fields)
			}
		case reflect.Struct:
			recursionField, err = d.decode(fv, src, fields)
		case reflect.Slice, reflect.Array:
			slice := reflect.MakeSlice(fv.Type(), len(src[name]), len(src[name]))
			for j, s := range src[name] {
				if err = d.decodeElement(fv.Type().Elem(), slice.Index(j), s); err != nil {
					goto End
				}
			}
			fv.Set(slice)
		case reflect.Map:
			mapField = fv
		default:
			if err = d.unmarshal(fv.Type(), fv, src.Get(name)); err != nil {
				goto End
			}
		}
	}

End:
	if mapField.IsValid() {
		return mapField, err
	}
	return recursionField, err
}

func (d *Decoder) unmarshal(t reflect.Type, v reflect.Value, src string) (err error) {
	switch t.Kind() {
	case reflect.Bool:
		val := reflect.New(BoolType)
		if err = val.Interface().(Unmarshaler).UnmarshalURL(src); err != nil {
			return
		}
		v.SetBool(val.Elem().Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := reflect.New(Int64Type)
		if err = val.Interface().(Unmarshaler).UnmarshalURL(src); err != nil {
			return
		}
		v.SetInt(val.Elem().Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := reflect.New(UInt64Type)
		if err = val.Interface().(Unmarshaler).UnmarshalURL(src); err != nil {
			return
		}
		v.SetUint(val.Elem().Uint())
	case reflect.Float32:
		val := reflect.New(Float32Type)
		if err = val.Interface().(Unmarshaler).UnmarshalURL(src); err != nil {
			return
		}
		v.SetFloat(val.Elem().Float())
	case reflect.Float64:
		val := reflect.New(Float64Type)
		if err = val.Interface().(Unmarshaler).UnmarshalURL(src); err != nil {
			return
		}
		v.SetFloat(val.Elem().Float())
	case reflect.String:
		val := reflect.New(StringType)
		_ = val.Interface().(Unmarshaler).UnmarshalURL(src) // Never return errors
		v.SetString(val.Elem().String())
	case reflect.Interface:
		val := reflect.New(InterfaceType)
		_ = val.Interface().(Unmarshaler).UnmarshalURL(src) // Never return errors
		v.Set(val.Elem().Field(0))
	//case reflect.Ptr:
	//case reflect.Slice, reflect.Array:
	//case reflect.Map:
	//case reflect.Struct:
	//case reflect.Complex64, reflect.Complex128:
	default:
		return fmt.Errorf("unmarshaler not found for %v", t)
	}

	return
}
