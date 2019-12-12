package form

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestBaseType(t *testing.T) {
	type TestType struct {
		BV bool `form:"b_v,omitempty"`
		UV uint
		SV []int          `form:"s_v"`
		MV map[int]string `form:"m_v"`
		EX []string       `form:"-"`
	}

	exp := url.Values{
		"11":  []string{"a"},
		"22":  []string{"b"},
		"UV":  []string{"12345"},
		"s_v": []string{"1", "2"},
	}
	v1 := TestType{
		SV: []int{1, 2},
		UV: 12345,
		MV: map[int]string{
			11: "a",
			22: "b",
		},
	}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
}

func TestFloatType(t *testing.T) {
	type TestType struct {
		FV32 float32 `form:"f_32"`
		FV64 float64
	}

	exp := url.Values{
		"f_32": []string{"2.7182817459106445"},
		"FV64": []string{"3.14159265358979311599796346854419"},
	}
	v1 := TestType{
		FV32: math.E,
		FV64: math.Pi,
	}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(float64(v1.FV32-v2.FV32)) > 0.0000000001 || math.Abs(v1.FV64-v2.FV64) > 0.0000000001 {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
}

func TestInterfaceType(t *testing.T) {
	type TestType struct {
		BV bool          `form:"b_v"`
		SV []interface{} `form:"s_v"`
	}

	exp := url.Values{
		"b_v": []string{"false"},
		"s_v": []string{"1", "b"},
	}
	v1 := TestType{
		SV: []interface{}{1, "b"}, // interface{} type will be decoded as string
	}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if v1.BV != v2.BV || len(v1.SV) != len(v2.SV) {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
	for i, v := range v1.SV {
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", v2.SV[i]) {
			t.Fatal("invalid decode result:", v2, "expected:", v1)
		}
	}
}

func TestEmptyEmbedType(t *testing.T) {
	type EmbedType struct {
	}
	type TestType struct {
		EmbedType
	}

	exp := url.Values{}
	v1 := TestType{}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
}

func TestOmitEmptyEmbedType(t *testing.T) {
	type EmbedType struct {
		V bool
	}
	type TestType struct {
		EmbedType `form:",omitempty"`
	}

	exp := url.Values{}
	v1 := TestType{
		EmbedType: EmbedType{V: false},
	}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
}

func TestNullPtrType(t *testing.T) {
	type EmbedType struct {
	}
	type TestType struct {
		*EmbedType
	}

	exp := url.Values{
		"EmbedType": []string{"null"},
	}
	v1 := TestType{}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
}

func TestPtrType(t *testing.T) {
	type EmbedType struct {
		IV int
	}
	type TestType struct {
		Embed *EmbedType
	}

	exp := url.Values{
		"IV": []string{"10"},
	}
	v1 := TestType{
		Embed: &EmbedType{IV: 10},
	}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
}

type CustomErrType struct {
	V int
}

func (e CustomErrType) MarshalURL() (string, error) {
	return "", errors.New("err")
}

func (e *CustomErrType) UnmarshalURL(v string) error {
	return errors.New("err")
}

type CustomBool bool

func (b CustomBool) String() string {
	if b {
		return "Y"
	}
	return "N"
}

func (b *CustomBool) MarshalURL() (string, error) {
	return b.String(), nil
}

func (b *CustomBool) UnmarshalURL(v string) error {
	if v == "Y" {
		*b = true
	}
	return nil
}

type CustomTime time.Time

func (t CustomTime) String() string {
	return time.Time(t).Format("20060102150405")
}

func (t *CustomTime) MarshalURL() (string, error) {
	return t.String(), nil
}

func (t *CustomTime) UnmarshalURL(v string) error {
	dt, err := time.ParseInLocation("20060102150405", v, time.Now().Location())
	if err != nil {
		return err
	}
	*t = CustomTime(dt)
	return nil
}

func TestCustomMarshalerType(t *testing.T) {
	type TestType struct {
		ATime *CustomTime
		BTime CustomTime
	}

	now := time.Now()
	aTime := CustomTime(now)

	v1 := TestType{
		ATime: &aTime,
		BTime: CustomTime(now.Add(time.Hour)),
	}
	exp := url.Values{
		"ATime": []string{v1.ATime.String()},
		"BTime": []string{v1.BTime.String()},
	}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if time.Time(*v1.ATime).Unix() != time.Time(*v2.ATime).Unix() ||
		time.Time(v1.BTime).Unix() != time.Time(v2.BTime).Unix() {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
}

func TestCustomMarshalerType2(t *testing.T) {
	type TestType struct {
		AVal CustomBool
		BVal *CustomBool
	}

	b := false
	v1 := TestType{
		AVal: CustomBool(true),
		BVal: (*CustomBool)(&b),
	}
	exp := url.Values{
		"AVal": []string{"Y"},
		"BVal": []string{"N"},
	}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
}

func TestCustomMarshalerType3(t *testing.T) {
	type TestType struct {
		SV []CustomBool
	}

	v1 := TestType{
		SV: []CustomBool{CustomBool(true), CustomBool(false)},
	}
	exp := url.Values{
		"SV": []string{"Y", "N"},
	}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
}

func TestCustomMarshalerType4(t *testing.T) {
	type TestType struct {
		EV CustomErrType
	}

	v := TestType{}

	// Encode
	if _, err := Encode(&v); err == nil {
		t.Fatal("expected err")
	}

	// Decode
	if err := Decode(&v, url.Values{}); err == nil {
		t.Fatal("expected err")
	}
}

func TestCustomMarshalerType5(t *testing.T) {
	type TestType struct {
		EV []*CustomErrType
	}

	v := TestType{
		EV: []*CustomErrType{{1}},
	}

	// Encode
	if _, err := Encode(&v); err == nil {
		t.Fatal("expected err")
	}

	// Decode
	if err := Decode(&v, url.Values{"EV": []string{"1"}}); err == nil {
		t.Fatal("expected err")
	}
}

func TestCustomMarshalerType6(t *testing.T) {
	type TestType struct {
		EV map[string]CustomErrType
	}

	v := TestType{
		EV: map[string]CustomErrType{
			"a": {1},
		},
	}
	if _, err := Encode(&v); err == nil {
		t.Fatal("expected err")
	}
}

func TestCustomMarshalerType7(t *testing.T) {
	type TestType struct {
		EV map[CustomErrType]int
	}

	v := TestType{
		EV: map[CustomErrType]int{
			{1}: 1,
		},
	}
	if _, err := Encode(&v); err == nil {
		t.Fatal("expected err")
	}
}

func TestCustomMarshalerType8(t *testing.T) {
	type Embed struct {
		EV *CustomErrType
	}
	type TestType struct {
		Embed
	}

	v := TestType{
		Embed: Embed{
			EV: &CustomErrType{1},
		},
	}
	if _, err := Encode(&v); err == nil {
		t.Fatal("expected err")
	}
}

func TestCustomMarshalerType9(t *testing.T) {
	type Embed struct {
		EV *CustomErrType
	}
	type TestType struct {
		*Embed
	}

	v := TestType{
		Embed: &Embed{
			EV: &CustomErrType{1},
		},
	}
	if _, err := Encode(&v); err == nil {
		t.Fatal("expected err")
	}
}

func TestEmptySlice(t *testing.T) {
	type TestType struct {
		SV []interface{} `form:"s_v,omitempty"`
	}

	exp := url.Values{}
	v1 := TestType{}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if len(v1.SV) != 0 || len(v2.SV) != 0 {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
}

func TestInvalidInputType(t *testing.T) {
	src := url.Values{}
	_, err := Encode(&src)
	if err != TypeError {
		t.Fatal("expected err:", TypeError, "returns:", err)
	}

	dst := map[string]string{}
	err = Decode(&dst, src)
	if err != TypeError {
		t.Fatal("expected err:", TypeError, "returns:", err)
	}
}

func TestInvalidDecode(t *testing.T) {
	type TestType struct {
		BV  bool
		IV  int
		UV  uint
		F32 float32
		F64 float64
	}
	src := url.Values{
		"BV":  []string{"true"},
		"IV":  []string{"-1"},
		"UV":  []string{"1"},
		"F32": []string{"3.14"},
		"F64": []string{"3.14"},
	}

	v := TestType{}
	// Bool
	src["BV"] = []string{""}
	if err := Decode(&v, src); err != nil {
		t.Fatal("excepted no error for bool", err)
	}
	src["BV"] = []string{"a"}
	if err := Decode(&v, src); err == nil {
		t.Fatal("excepted error for bool", v)
	}
	src["BV"] = []string{"true"}

	// Int
	src["IV"] = []string{""}
	if err := Decode(&v, src); err != nil {
		t.Fatal("excepted no error for int", err)
	}
	src["IV"] = []string{"a"}
	if err := Decode(&v, src); err == nil {
		t.Fatal("excepted error for int", v)
	}
	src["IV"] = []string{"-1"}

	// UInt
	src["UV"] = []string{}
	if err := Decode(&v, src); err != nil {
		t.Fatal("excepted no error for uint", err)
	}
	src["UV"] = []string{"a"}
	if err := Decode(&v, src); err == nil {
		t.Fatal("excepted error for uint", v)
	}
	src["UV"] = []string{"1"}

	// F32
	src["F32"] = []string{}
	if err := Decode(&v, src); err != nil {
		t.Fatal("excepted no error for float32", err)
	}
	src["F32"] = []string{"a"}
	if err := Decode(&v, src); err == nil {
		t.Fatal("excepted error for float32", v)
	}
	src["F32"] = []string{"3.14"}

	// F64
	src["F64"] = []string{}
	if err := Decode(&v, src); err != nil {
		t.Fatal("excepted no error for float64", err)
	}
	src["F64"] = []string{"a"}
	if err := Decode(&v, src); err == nil {
		t.Fatal("excepted error for float64", v)
	}
	src["F64"] = []string{"3.14"}
}

func TestComplex(t *testing.T) {
	type TestType struct {
		CV complex64
	}

	v := TestType{
		CV: complex(10, 10),
	}

	// Encode
	if _, err := Encode(&v); err == nil {
		t.Fatal("expected err")
	}

	// Decode
	if err := Decode(&v, url.Values{"CV": []string{}}); err == nil {
		t.Fatal("expected err")
	}
}

func TestMap(t *testing.T) {
	type TestType struct {
		M map[int]int
	}

	exp := url.Values{
		"1": []string{"111"},
		"2": []string{"222"},
	}
	v1 := TestType{
		M: map[int]int{
			1: 111,
			2: 222,
		},
	}

	// Encode
	val, err := Encode(&v1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, exp) {
		t.Fatal("invalid encode result:", val, "expected:", exp)
	}

	// Decode
	exp["3"] = []string{"ccc"}
	exp["a"] = []string{"aaa"}
	v1.M[3] = 0

	v2 := TestType{}
	err = Decode(&v2, exp)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v1, v2) {
		t.Fatal("invalid decode result:", v2, "expected:", v1)
	}
}
