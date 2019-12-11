package form

import (
	"net/url"
	"testing"
)

func TestDecoder_Decode(t *testing.T) {
	tt := Test{}
	vals, err := url.ParseQuery("IV=12&SliceVal=xx&SliceVal=aaa&Time=20191211150743&b_val=true&m1=1&m2=2")
	if err != nil {
		t.Fatal(err)
	}
	err = NewDecoder().Decode(&tt, vals)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tt)
}
