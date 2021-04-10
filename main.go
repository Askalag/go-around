package main

import (
	"context"
	"github.com/Askalag/go-around/handlers"
	"github.com/go-openapi/runtime/middleware"
	gHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main()  {
	l := log.New(os.Stdout, "go-around ", log.LstdFlags)
	//hh := handlers.NewHello(l)
	//bye := handlers.NewGoodbye(l)

	const basePath = "./images"
	const swaggerFilePath = "/api-docs/swagger.yaml"

	// ReDoc and swagger config
	opts := middleware.RedocOpts{SpecURL: swaggerFilePath}
	sh := middleware.Redoc(opts, nil)

	p := handlers.NewProduct(l)
	//f := handlers.NewFiles(l, basePath)
	router := mux.NewRouter()

	// simple crud
	router.HandleFunc("/", p.GetProducts).Methods("GET")
	router.HandleFunc("/{id:[0-9]+}", p.UpdateProduct).Methods("PUT")
	router.HandleFunc("/", p.AddProduct).Methods("POST")
	// todo delete method
	//router.Use(p.MiddlewareProductValidation)

	// file
	router.HandleFunc("/images/{id}", func (rw http.ResponseWriter, req *http.Request) {})

	// swagger and reDoc
	router.Handle("/docs", sh).Methods("GET")
	router.Handle(swaggerFilePath, http.FileServer(http.Dir("./"))).Methods("GET")

	//router.Handle("/products", p)
	//mux.Handle("/bye", bye)
	//mux.Handle("/pl", pl)


	// gzip
	gh := gHandlers.CompressHandler(router)

	// CORS for angular
	ch := gHandlers.CORS(gHandlers.AllowedOrigins([]string{"http://localhost:4200"}))
	ch(gh)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      ch(gh),
		ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}


	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)


	sig := <- sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30 * time.Second)
	_ = s.Shutdown(tc)


}
