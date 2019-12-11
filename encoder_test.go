package form

import (
	"net/url"
	"testing"
	"time"
)

type Embed struct {
	IV int
	SV string `form:",omitempty"`
}

type Test struct {
	BValue    bool `form:"b_val"`
	IntVal    int   `form:"int_val,omitempty"`
	SliceVal  []interface{}
	MapVal    map[string]int `form:",omitempty"`
	StructVal struct{}
	Time      TTime
	EM  Embed
}

type TTime time.Time

func (t TTime) MarshalURL() (string, error) {
	return time.Time(t).Format("20060102150405"), nil
}

func (t *TTime) UnmarshalURL(v string) error {
	dt, err := time.Parse("20060102150405", v)
	if err != nil {
		return err
	}
	*t = TTime(dt)
	return nil
}

func TestEncoder_Encode(t *testing.T) {
	vals := url.Values{}
	encoder := NewEncoder()
	err := encoder.Encode(&Test{
		SliceVal: []interface{}{"xx", "aaa"},
		MapVal:   map[string]int{"m1": 1, "m2": 2},
		//Time:     TTime(time.Now()),
		//EM: Embed{},
	}, vals)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vals, vals.Encode())
}
