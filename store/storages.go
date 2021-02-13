package store

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"regexp"
	"time"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", ValidateSKU)
	return validate.Struct(p)
}

func ValidateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-absd-df
	fieldValue := fl.Field().String()
	regx := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := regx.FindAllString(fieldValue, -1)

	if len(matches) != 1 {
		return false
	}
	return true
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