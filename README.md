# go-form

`go-form` converts a struct and `url.Values` to each other which inspired by [`gorilla/schema`](https://github.com/gorilla/schema).

## Example

Here's a quick example: we parse POST form values and then decode them into a struct:

```go
type Person struct {
    Name  string
    Phone string
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        // Handle error
    }

    var person Person

    // r.PostForm is a map of our POST form values
    err = form.Unmarshal(&person, r.PostForm)
    if err != nil {
        // Handle error
    }

    // Do something with person.Name or person.Phone
}
```

Conversely, contents of a struct can be encoded into form values. Here's a variant of the previous example:

```go
func MyHttpRequest() {
    person := Person{"Jane Doe", "555-5555"}

    vals, err := form.Marshal(&person)

    if err != nil {
        // Handle error
    }

    // Use form values, for example, with an http client
    client := new(http.Client)
    res, err := client.PostForm("http://my-api.test", vals)
}
```

To define custom names for fields, use a struct tag "form". To not populate certain fields, use a dash for the name and it will be ignored:

```go
type Person struct {
    Name  string `form:"name"`  // custom name
    Phone string `form:"phone"` // custom name
    Admin bool   `form:"-"`     // this field is never set
}
```

The supported field types in the struct are:

* bool
* float variants (float32, float64)
* int variants (int, int8, int16, int32, int64)
* string
* uint variants (uint, uint8, uint16, uint32, uint64)
* struct
* a pointer to one of the above types
* a slice of one of the above types or interface{} type
* a map of any above types
* custom types implements Marshaler and Unmarshaler interfaces

## Custom type implementation

```go
type Bool bool

func (b Bool) String() string {
	if b {
		return "Y"
	}
	return "N"
}

func (b Bool) MarshalURL() (string, error) {
	return b.String(), nil
}

func (b *Bool) UnmarshalURL(v string) error {
	if v == "Y" {
		*b = true
	}
	return nil
}
```

## Thanks

* [gorilla/schema](https://github.com/gorilla/schema)
