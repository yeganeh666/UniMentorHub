package routers

import (
	"Atrovan_Q1/controllers"
	"Atrovan_Q1/core/authentication"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// SetUsersRoutes include all routes belongs to users.
func SetUsersRoutes(router *mux.Router) *mux.Router {
	router.Handle("/lessons",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.GetAllLessonsCache),
			negroni.HandlerFunc(controllers.GetAllLessons),
		)).Methods("GET")

	router.Handle("/users/lessons",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.GetUserAllLessons),
		)).Methods("GET")

	router.Handle("/users/student/lesson",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RoleBasedAuthentication("student").ServeHTTP),
			negroni.HandlerFunc(controllers.SignUpForLesson),
		)).Methods("POST")

	router.Handle("/users/professor/lesson",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RoleBasedAuthentication("professor").ServeHTTP),
			negroni.HandlerFunc(controllers.CreateLesson),
		)).Methods("POST")

	return router
}
