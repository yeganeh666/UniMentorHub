package tests

import (
	"Atrovan_Q1/core/authentication"
	"Atrovan_Q1/db"
	"Atrovan_Q1/services/models"
	"os"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	os.Setenv("JWT_PRIVATE_KEY_PATH", "../config/keys/private_key.pem")
	os.Setenv("JWT_PUBLIC_KEY_PATH", "../config/keys/public_key.pub")
}

func TestInitJWTAuthenticationBackend(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	assert.NotNil(t, authBackend)
	assert.NotNil(t, authBackend.PublicKey)
}

func TestGenerateToken(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken(1, "student")

	assert.Nil(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})

	assert.Nil(t, err)
	assert.True(t, token.Valid)
}

func TestAuthenticate(t *testing.T) {
	var userModel db.UserModel
	authBackend := authentication.InitJWTAuthenticationBackend()
	user := &models.User{
		UserID:   11111111,
		Password: "test",
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hashedPassword)
	userModel.InsertOne(user)

	user.Password = "test"

	ok, _ := authBackend.Authenticate(user)
	assert.Equal(t, ok, true)
}

func TestAuthenticateIncorrectPass(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	user := &models.User{
		UserID:   11111111,
		Password: "test-pass",
	}
	ok, role := authBackend.Authenticate(user)

	assert.Equal(t, role, "")
	assert.Equal(t, ok, false)
}

func TestAuthenticateIncorrectUserID(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	user := &models.User{
		UserID:   123,
		Password: "test",
	}
	ok, role := authBackend.Authenticate(user)

	assert.Equal(t, role, "")
	assert.Equal(t, ok, false)
}

func TestLogoutnJWT(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken(1, "student")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	err = authBackend.Logout(tokenString, token)
	assert.Nil(t, err)

	redisConn, _ := db.ConnectRedis()
	redisValue, err := redisConn.GetValue(tokenString)
	assert.Nil(t, err)
	assert.NotEmpty(t, redisValue)
}

func TestIsInBlacklist(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken(1, "student")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	err = authBackend.Logout(tokenString, token)
	assert.Nil(t, err)

	assert.True(t, authBackend.IsInBlacklist(tokenString))
}

func TestIsNotInBlacklist(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	assert.False(t, authBackend.IsInBlacklist("1"))
}
