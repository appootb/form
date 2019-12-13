package form

import (
	"errors"
	"net/url"
)

var (
	TypeError = errors.New("the interface must be a pointer to a struct")
)

func Marshal(src interface{}) (url.Values, error) {
	v := url.Values{}
	err := NewEncoder(v).Encode(src)
	return v, err
}

func Unmarshal(dst interface{}, src url.Values) error {
	return NewDecoder(src).Decode(dst)
}
