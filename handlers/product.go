package handlers

import (
	"github.com/Askalag/go-around/store"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request)  {

	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		regx := regexp.MustCompile(`/([0-9]+)`)
		g := regx.FindAllStringSubmatch(r.URL.Path, 2)

		if len(g) != 1 {
			http.Error(rw, "Invalid Url", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(rw, "Invalid Url", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, _ := strconv.Atoi(idString)

		p.updateProduct(id, rw, r)
		return
	}
	// catch all ...
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT add product request")
	prod := &store.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
		return
	}
	p.l.Printf("Prod: %#v", prod)

	res := store.UpdateProduct(id, prod)
	if res != nil {
		http.Error(rw, res.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST add product request")
	prod := &store.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}
	p.l.Printf("Prod: %#v", prod)
	store.AddProduct(prod)

}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request)  {
	lp := store.GetStaticProducts()
	//rw.Header().Add("Content-Type", "application/json")
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}