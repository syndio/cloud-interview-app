package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	employeesProxy, err := newProxy("http://employeesapi:6541")
	if err != nil {
		log.Panic(err)
	}
	r.Handle("/employees*", employeesProxy)

	log.Printf("listening on %s ...", os.Getenv("PORT"))
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r))
}

func newProxy(targetURL string) (http.Handler, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return httputil.NewSingleHostReverseProxy(u), nil
}
