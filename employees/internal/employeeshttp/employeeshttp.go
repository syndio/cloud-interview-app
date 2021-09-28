package employeeshttp

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/syndio/cloud-interview-app/employees/internal/employeesdb"
)

var cacheKey = "employees"

type Handler struct {
	db    *employeesdb.Database
	cache redis.UniversalClient
}

func NewHandler(db *employeesdb.Database, cache redis.UniversalClient) *Handler {
	return &Handler{db: db, cache: cache}
}

func (h *Handler) fromCache(ctx context.Context, key string, result interface{}) (bool, error) {
	res, err := h.cache.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error getting from cache: %w", err)
	}
	if err := json.Unmarshal([]byte(res), result); err != nil {
		return true, fmt.Errorf("error unmarshaling result: %w", err)
	}
	return true, nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var employee *employeesdb.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := h.db.Create(ctx, employee); err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.cache.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(employee)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		log.Printf("%+v", err)
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var employees []*employeesdb.Employee
	isCached, err := h.fromCache(ctx, cacheKey, &employees)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isCached {
		employees, err = h.db.List(ctx)
		if err != nil {
			log.Printf("%+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	b, err := json.Marshal(employees)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.cache.Set(ctx, cacheKey, b, 0).Err(); err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		log.Printf("%+v", err)
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	employeeID, err := strconv.ParseInt(chi.URLParam(r, "employeeID"), 10, 64)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.db.Delete(ctx, employeeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.cache.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
