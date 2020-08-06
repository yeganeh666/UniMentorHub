package main

import (
	"Atrovan_Q1/db"
	"Atrovan_Q1/routers"
	"Atrovan_Q1/validators"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var router *mux.Router

func init() {
	validators.Init()
	db.InitConnections()
	router = routers.InitRoutes()
}

func main() {
	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServe(":5000", n)
}
