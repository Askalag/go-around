// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"github.com/Askalag/go-around/store"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT add product request")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	prod := r.Context().Value(KeyProduct{}).(store.Product)
	p.l.Printf("Prod: %#v", prod)

	res := store.UpdateProduct(id, &prod)
	if res != nil {
		http.Error(rw, res.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

// swagger:route POST /products products product
// adding another product
// responses:
// 201:

// Add a new product
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST add product request")
	prod := r.Context().Value(KeyProduct{}).(*store.Product)
	store.AddProduct(prod)
	rw.WriteHeader(http.StatusOK)
}

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 200: productResponse

// GetProducts returns the products from the static data
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request)  {
	lp := store.GetStaticProducts()
	//rw.Header().Add("Content-Type", "application/json")
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

type KeyProduct struct {}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func (rw http.ResponseWriter, r *http.Request) {
		prod := &store.Product{}
		err := prod.FromJSON(r.Body)

		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
			return
		}

		err2 := prod.Validate()

		if err2 != nil {
			http.Error(rw, "Bad Values for Product", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}