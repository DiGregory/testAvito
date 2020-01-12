package main

import (
	"github.com/go-chi/chi"
	"net/http"
	"fmt"
)

func (a *apiApp) createApiHandlers() (error) {
	r := chi.NewRouter()
	r.Get("/adverts", a.getAdvertsHandler)
	r.Get("/adverts/{id}", a.getSingleAdvertHandler)
	r.Post("/adverts", a.createAdvertHandler)

	fmt.Println("Server started at ", a.Addr)
	return http.ListenAndServe(a.Addr, r)
}
