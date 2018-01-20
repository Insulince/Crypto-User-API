package models

import "github.com/gorilla/mux"

type Router struct {
	*mux.Router
}

func CreateRouter() (router *Router) {
	return &Router{
		mux.NewRouter().StrictSlash(true),
	}
}
