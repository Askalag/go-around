package store

import "testing"

func TestProductValidation(t *testing.T) {
	p := &Product{}
	p.Name = "okok"
	p.Price = 8.2
	p.SKU = "rfrfrfrfe-l-l"

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}