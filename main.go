package main

import (
	"context"
	"github.com/Askalag/go-around/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main()  {
	l := log.New(os.Stdout, "go-around ", log.LstdFlags)
	hh := handlers.NewHello(l)
	bye := handlers.NewGoodbye(l)
	pl := handlers.NewProduct(l)

	mux := http.NewServeMux()
	mux.Handle("/", hh)
	mux.Handle("/bye", bye)
	mux.Handle("/pl", pl)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      mux,
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
