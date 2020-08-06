package authentication

import (
	"Atrovan_Q1/services/models"
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
)

// RoleBasedMiddleware keeps the minimum role that route need to it.
type RoleBasedMiddleware struct {
	role string
}

// RequireTokenAuthentication is a middleware to check if the user have token or not.
func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authBackend := InitJWTAuthenticationBackend()

	token, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return authBackend.PublicKey, nil
	})

	if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
		req.Header.Set("id", fmt.Sprint(token.Claims.(jwt.MapClaims)["sub"]))
		req.Header.Set("role", token.Claims.(jwt.MapClaims)["role"].(string))
		next(rw, req)
	} else {
		response, _ := json.Marshal(models.ErrorResponse{
			Code:      "Unauthorized",
			Message:   "Your request has not been applied because it lacks valid authentication credentials",
			MessageFa: "درخواست شما شامل کد امنیتی معتبر نمیباشد",
		})
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(response)
	}
}

// RoleBasedAuthentication check user limitions on doing operations.
func RoleBasedAuthentication(role string) *RoleBasedMiddleware {
	return &RoleBasedMiddleware{
		role: role,
	}
}

func (r *RoleBasedMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	role := req.Header.Get("role")
	if role == r.role {
		next(rw, req)
	} else {
		response, _ := json.Marshal(models.ErrorResponse{
			Code:      "Unauthorized",
			Message:   "Your request has not been applied because it lacks valid authentication credentials",
			MessageFa: "درخواست شما شامل کد امنیتی معتبر نمیباشد",
		})
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(response)
	}
}
