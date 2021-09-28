package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/syndio/cloud-interview-app/employees/internal/employeesdb"
	"github.com/syndio/cloud-interview-app/employees/internal/employeeshttp"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db, err := sql.Open("postgres", "host=postgres user=dev password=dev dbname=employees sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	cache := redis.NewClient(&redis.Options{Addr: "redis:6379"})

	eh := employeeshttp.NewHandler(employeesdb.NewDatabase(db), cache)
	r.Route("/employees", func(r chi.Router) {
		r.Get("/", eh.List)
		r.Post("/", eh.Create)
		r.Route("/{employeeID}", func(r chi.Router) {
			r.Delete("/", eh.Delete)
		})
	})

	log.Printf("listening on %s ...", os.Getenv("PORT"))
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r))
}
