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
	p := handlers.NewProduct(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", p.GetProducts)

	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", p.UpdateProduct)
	putRouter.Use(p.MiddlewareProductValidation)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", p.AddProduct)
	postRouter.Use(p.MiddlewareProductValidation)

	// ReDoc config
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
 	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	//sm.Handle("/products", p)
	//mux.Handle("/bye", bye)
	//mux.Handle("/pl", pl)


	// gzip
	gh := gHandlers.CompressHandler(sm)

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
