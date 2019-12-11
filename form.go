package form

import "net/url"

func Encode(src interface{}) (url.Values, error) {
	v := url.Values{}
	err := NewEncoder().Encode(src, v)
	return v, err
}

func Decode(dst interface{}, src url.Values) error {
	return NewDecoder().Decode(dst, src)
}
