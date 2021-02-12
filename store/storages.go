package store

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func GetStaticProducts() Products {
	return productList
}

func UpdateProduct(id int, p *Product) error {
	fp, i, err := findProduct(id)
	if err != nil {
		return err
	}
	p.Id = fp.Id
	productList[i] = p
	return nil
}

func findProduct(id int) (*Product, int,  error) {
	for i, prod := range productList {
		if prod.Id == id {
			return prod, i, nil
		}
	}
	return nil, -1, fmt.Errorf("product with id : %v not found", id)
}

func AddProduct(p *Product)  {
	p.Id = getNextId()
	productList = append(productList, p)
}

func getNextId() int {
	lp := productList[len(productList) -1]
	return lp.Id + 1
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

var productList = []*Product {
	{
		Id:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		Id:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "abc34da",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}