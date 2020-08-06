package routers

import (
	"github.com/gorilla/mux"
)

// InitRoutes gather all the routes together return them as router.
func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = SetUsersRoutes(router)
	router = SetAuthenticationRoutes(router)
	return router
}
