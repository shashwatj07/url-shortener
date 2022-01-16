package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/basic"
	"github.com/shaj13/go-guardian/auth/strategies/bearer"
	"github.com/shaj13/go-guardian/store"
)

// The Authenticator object.
var authenticator auth.Authenticator

// Cache to store the Authorzation Tokens.
var cache store.Cache

// Creates a JWT complaint Authorzation Token.
func CreateToken(c *gin.Context) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth-app",
		"sub": "cs559",
		"aud": "any",
		"exp": time.Now().AddDate(0, 0, 30).Unix(),
	})
	jwtToken, _ := token.SignedString([]byte("secret"))
	c.Writer.Write([]byte(jwtToken))
}

// Validates the user to check if his credentials are registered
// in our application.
func ValidateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	if userName == "cs559" && password == "iitbh" {
		return auth.NewDefaultUser("cs559", "1", nil, nil), nil
	}
	return nil, fmt.Errorf("invalid credentials")
}

// Verifies the Authorzation Token to determine if the request is
// to be allowed or not.
func VerifyToken(ctx context.Context, r *http.Request, tokenString string) (auth.Info, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user *auth.DefaultUser = auth.NewDefaultUser(claims["sub"].(string), "", nil, nil)
		return user, nil
	}
	return nil, fmt.Errorf("invaled token")
}

// Initializes a Go Guardian authentication instance to start
// a scalable authentication management system and assign the
// necessary strategies.
func SetupGoGuardian() {
	authenticator = auth.New()
	cache = store.NewFIFO(context.Background(), time.Minute*5)
	basicStrategy := basic.New(ValidateUser, cache)
	tokenStrategy := bearer.New(VerifyToken, cache)
	authenticator.EnableStrategy(basic.StrategyKey, basicStrategy)
	authenticator.EnableStrategy(bearer.CachedStrategyKey, tokenStrategy)
}
