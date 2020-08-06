package routers

import (
	"Atrovan_Q1/controllers"
	"Atrovan_Q1/core/authentication"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// SetAuthenticationRoutes include all routes belongs to authentication.
func SetAuthenticationRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/register", controllers.CreateNewUser).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.Handle("/refresh-token",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.RefreshToken),
		)).Methods("GET")
	router.Handle("/logout",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.Logout),
		)).Methods("GET")
	return router
}
