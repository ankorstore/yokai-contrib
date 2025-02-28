package model

import "github.com/google/jsonapi"

type Foo struct {
	ID   int    `jsonapi:"primary,foo"`
	Name string `jsonapi:"attr,name"`
	Bar  *Bar   `jsonapi:"relation,bar"`
}

func (f Foo) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"meta": "foo",
	}
}

type Bar struct {
	ID   int    `jsonapi:"primary,bar"`
	Name string `jsonapi:"attr,name"`
}

func (b Bar) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"meta": "bar",
	}
}

func CreateTestFoo() Foo {
	return Foo{
		ID:   123,
		Name: "foo",
		Bar: &Bar{
			ID:   456,
			Name: "bar",
		},
	}
}
