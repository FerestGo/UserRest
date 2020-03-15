package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getRouter() *mux.Router {
	router := mux.NewRouter()

	userApi := UserApi{}
	router.HandleFunc(userApi.PathWithId(), userApi.One).Methods(http.MethodGet)
	router.HandleFunc(userApi.Path(), userApi.Many).Methods(http.MethodGet)
	router.HandleFunc(userApi.Path(), userApi.Create).Methods(http.MethodPost)

	router.HandleFunc(userApi.PathWithId(), userApi.Update).Methods(http.MethodPut)
	return router
}
