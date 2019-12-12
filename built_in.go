package form

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	NullValue = "null"
)

var (
	BoolType      = reflect.TypeOf(Bool(false))
	Int64Type     = reflect.TypeOf(Int64(0))
	UInt64Type    = reflect.TypeOf(UInt64(0))
	Float32Type   = reflect.TypeOf(Float32(0.0))
	Float64Type   = reflect.TypeOf(Float64(0.0))
	StringType    = reflect.TypeOf(String(""))
	InterfaceType = reflect.TypeOf(Interface{})
)

type Bool bool

func (v Bool) MarshalURL() (string, error) {
	return strconv.FormatBool(bool(v)), nil
}

func (v *Bool) UnmarshalURL(src string) error {
	if src == "" {
		return nil
	}
	val, err := strconv.ParseBool(src)
	if err != nil {
		return err
	}
	*v = Bool(val)
	return nil
}

type Int64 int64

func (v Int64) MarshalURL() (string, error) {
	return strconv.FormatInt(int64(v), 10), nil
}

func (v *Int64) UnmarshalURL(src string) error {
	if src == "" {
		return nil
	}
	val, err := strconv.ParseInt(src, 10, 0)
	if err != nil {
		return err
	}
	*v = Int64(val)
	return nil
}

type UInt64 uint64

func (v UInt64) MarshalURL() (string, error) {
	return strconv.FormatUint(uint64(v), 10), nil
}

func (v *UInt64) UnmarshalURL(src string) error {
	if src == "" {
		return nil
	}
	val, err := strconv.ParseUint(src, 10, 0)
	if err != nil {
		return err
	}
	*v = UInt64(val)
	return nil
}

type Float32 float32

func (v Float32) MarshalURL() (string, error) {
	return strconv.FormatFloat(float64(v), 'f', 16, 32), nil
}

func (v *Float32) UnmarshalURL(src string) error {
	if src == "" {
		return nil
	}
	val, err := strconv.ParseFloat(src, 32)
	if err != nil {
		return err
	}
	*v = Float32(val)
	return nil
}

type Float64 float64

func (v Float64) MarshalURL() (string, error) {
	return strconv.FormatFloat(float64(v), 'f', 32, 64), nil
}

func (v *Float64) UnmarshalURL(src string) error {
	if src == "" {
		return nil
	}
	val, err := strconv.ParseFloat(src, 64)
	if err != nil {
		return err
	}
	*v = Float64(val)
	return nil
}

type String string

func (v String) MarshalURL() (string, error) {
	return string(v), nil
}

func (v *String) UnmarshalURL(src string) error {
	*v = String(src)
	return nil
}

type Interface struct {
	Val interface{}
}

func (v Interface) MarshalURL() (string, error) {
	return fmt.Sprintf("%v", v.Val), nil
}

func (v *Interface) UnmarshalURL(src string) error {
	v.Val = src
	return nil
}
