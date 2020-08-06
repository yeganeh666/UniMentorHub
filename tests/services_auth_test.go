package tests

import (
	"Atrovan_Q1/core/authentication"
	"Atrovan_Q1/services"
	"Atrovan_Q1/services/models"

	"net/http"
	"os"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("JWT_PRIVATE_KEY_PATH", "../config/keys/private_key.pem")
	os.Setenv("JWT_PUBLIC_KEY_PATH", "../config/keys/public_key.pub")
}

func TestLoginIncorrectPassword(t *testing.T) {
	user := models.User{
		UserID:   11111111,
		Password: "test-pass",
	}
	response, token := services.Login(&user)
	assert.Equal(t, http.StatusUnauthorized, response)
	assert.Empty(t, token)
}

func TestLoginIncorrectUserID(t *testing.T) {
	user := models.User{
		UserID:   123,
		Password: "test",
	}
	response, token := services.Login(&user)
	assert.Equal(t, http.StatusUnauthorized, response)
	assert.Empty(t, token)
}

func TestLoginEmptyCredentials(t *testing.T) {
	user := models.User{
		Password: "",
	}
	response, token := services.Login(&user)
	assert.Equal(t, http.StatusUnauthorized, response)
	assert.Empty(t, token)
}

func TestRefreshToken(t *testing.T) {
	user := models.User{
		UserID:   11111111,
		Password: "test",
	}
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken(user.UserID, "student")
	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	assert.Nil(t, err)

	newToken := services.RefreshToken(&user, "")
	assert.NotEmpty(t, newToken)
}
