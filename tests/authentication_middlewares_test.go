package tests

import (
	"Atrovan_Q1/core/authentication"
	"Atrovan_Q1/routers"
	"Atrovan_Q1/services"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/stretchr/testify/assert"
)

var (
	token       string
	server      *negroni.Negroni
	authBackend *authentication.JWTAuthenticationBackend
)

func init() {
	os.Setenv("JWT_PRIVATE_KEY_PATH", "../config/keys/private_key.pem")
	os.Setenv("JWT_PUBLIC_KEY_PATH", "../config/keys/public_key.pub")
	authBackend = authentication.InitJWTAuthenticationBackend()
	token, _ = authBackend.GenerateToken(1, "student")

	router := routers.InitRoutes()
	server = negroni.Classic()
	server.UseHandler(router)
}

func TestRequireTokenAuthentication(t *testing.T) {
	resource := "/lessons"

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", resource, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
}

func TestRequireTokenAuthenticationInvalidToken(t *testing.T) {
	resource := "/lessons"

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", resource, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", "token"))
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusUnauthorized)
}

func TestRequireTokenAuthenticationEmptyToken(t *testing.T) {
	resource := "/lessons"

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", resource, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ""))
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusUnauthorized)
}

func TestRequireTokenAuthenticationWithoutToken(t *testing.T) {
	resource := "/lessons"

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", resource, nil)
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusUnauthorized)
}

func TestRoleBasedAuthenticationWithTokenProfessorForProfessorRouteWithEmptyBody(t *testing.T) {
	resource := "/professor/lesson"
	token, _ = authBackend.GenerateToken(1, "professor")

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", resource, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusBadRequest)
}

func TestRoleBasedAuthenticationWithTokenStudentForProfessorRoute(t *testing.T) {
	resource := "/professor/lesson"
	token, _ = authBackend.GenerateToken(1, "student")

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", resource, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusUnauthorized)
}

func TestRequireTokenAuthenticationAfterLogout(t *testing.T) {
	resource := "/lessons"

	requestLogout, _ := http.NewRequest("GET", resource, nil)
	requestLogout.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	services.Logout(requestLogout)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", resource, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusUnauthorized)
}
