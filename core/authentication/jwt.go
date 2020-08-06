package authentication

import (
	"Atrovan_Q1/config"
	"Atrovan_Q1/db"
	"Atrovan_Q1/services/models"
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// JWTAuthenticationBackend struct for do operations for authentication and holds keys.
type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// InitJWTAuthenticationBackend prepare the private and public key for sign the jwt.
func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}
	return authBackendInstance
}

// jwt default values.
var (
	jwtSecret     = config.EnvGetStr("JWT_SECRET", "")
	tokenDuration = config.EnvGetInt("JWT_TOKEN_DURATION", 25)
	expireOffset  = config.EnvGetInt("JWT_EXPIRE_OFFSET", 3600)
)

var authBackendInstance *JWTAuthenticationBackend = nil

// GenerateToken generate new jwt token based on userId.
func (backend *JWTAuthenticationBackend) GenerateToken(userID int, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims = jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * time.Duration(tokenDuration)).Unix(),
		"iat":  time.Now().Unix(),
		"sub":  userID,
		"role": role,
	}
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
	}
	return tokenString, nil
}

// Authenticate checks the hash password and returns true if user entered the correct credentials.
func (backend *JWTAuthenticationBackend) Authenticate(u *models.User) (bool, string) {
	var userModel db.UserModel
	user, err := userModel.FindOne(u)
	if err != nil {
		return false, ""
	}
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)) == nil, user.Role
}

// getTokenRemainingValidity calculate the remainer time of jwt token base on it's [exp].
func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) time.Duration {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return remainer + time.Second*time.Duration(expireOffset)
		}
	}
	return time.Second * time.Duration(expireOffset)
}

// Logout puts users jwt token to redis .
func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	redisConn, err := db.ConnectRedis()
	if err != nil {
		return fmt.Errorf("redis not connect")
	}
	return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(token.Claims.(jwt.MapClaims)["exp"]))
}

// IsInBlacklist check redis database to find if the jwt is valid or not
// when user logout from the app, his jwt is valid if the [exp] is bigger than time.now
// because of that we insert this jwt token to redis with remainer [exp] + offset ttl to redis
func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	redisConn, _ := db.ConnectRedis()
	redisToken, err := redisConn.GetValue(token)

	if err != nil {
		fmt.Println(err.Error())
	}

	if redisToken == "" {
		return false
	}

	return true
}

// getPrivateKey opens the PrivateKey file and returns it as *rsa.PrivateKey.
func getPrivateKey() *rsa.PrivateKey {
	absFilePath, _ := filepath.Abs(config.EnvGetStr("JWT_PRIVATE_KEY_PATH", "./config/keys/private_key.pem"))
	privateKeyFile, err := os.Open(absFilePath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))
	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

// getPublicKey opens the publicKey file and returns it as *rsa.publicKey.
func getPublicKey() *rsa.PublicKey {
	absFilePath, _ := filepath.Abs(config.EnvGetStr("JWT_PUBLIC_KEY_PATH", "./config/keys/public_key.pub"))
	publicKeyFile, err := os.Open(absFilePath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))
	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)
	if !ok {
		panic(err)
	}

	return rsaPub
}
