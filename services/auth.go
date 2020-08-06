package services

import (
	"Atrovan_Q1/core/authentication"
	"Atrovan_Q1/db"
	"Atrovan_Q1/services/models"
	"Atrovan_Q1/validators"
	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
	"golang.org/x/crypto/bcrypt"
)

// TokenAuthentication handles a type for jwt token.
type TokenAuthentication struct {
	Token string `json:"token" form:"token"`
}

// Login service
// check if authenticate status
// generate token based on user userId
func Login(user *models.User) (int, []byte) {
	authBackend := *authentication.InitJWTAuthenticationBackend()

	if ok, role := authBackend.Authenticate(user); ok {
		token, err := authBackend.GenerateToken(user.UserID, role)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		}
		response, _ := json.Marshal(TokenAuthentication{token})
		return http.StatusOK, response
	}

	return http.StatusUnauthorized, []byte("")
}

// RefreshToken service
// generate token based on user
func RefreshToken(user *models.User, role string) []byte {
	authBackend := *authentication.InitJWTAuthenticationBackend()
	token, err := authBackend.GenerateToken(user.UserID, role)
	if err != nil {
		panic(err)
	}
	response, err := json.Marshal(TokenAuthentication{token})
	if err != nil {
		panic(err)
	}
	return response
}

// Logout service
// get the token from request header
// use authBackend.Logout to handle logout from system
func Logout(req *http.Request) error {
	authBackend := *authentication.InitJWTAuthenticationBackend()
	tokenRequest, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		return err
	}
	tokenString := req.Header.Get("Authorization")
	return authBackend.Logout(tokenString, tokenRequest)
}

// CreateNewUser service
// validate json for type matching
// hash password
// validate based on our needs
// create user if not exists and returns token
func CreateNewUser(user *models.User) (int, []byte) {
	var userModel db.UserModel
	if ok, response := validators.Validate(user); !ok {
		return response.StatusCode, response.Body
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return models.ServerDownResponse.StatusCode, models.ServerDownResponse.Body
	}
	user.Password = string(hashedPassword)
	if ok, response := userModel.InsertOne(user); !ok {
		return response.StatusCode, response.Body
	}
	authBackend := *authentication.InitJWTAuthenticationBackend()
	token, err := authBackend.GenerateToken(user.UserID, user.Role)
	if err != nil {
		return models.ServerDownResponse.StatusCode, models.ServerDownResponse.Body
	}
	response, _ := json.Marshal(TokenAuthentication{token})
	return http.StatusOK, response
}
