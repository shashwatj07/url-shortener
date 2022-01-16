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

var authenticator auth.Authenticator
var cache store.Cache

func CreateToken(c *gin.Context) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth-app",
		"sub": "cs559",
		"aud": "any",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	jwtToken, _ := token.SignedString([]byte("secret"))
	c.Writer.Write([]byte(jwtToken))
}

func ValidateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	if userName == "cs559" && password == "iitbh" {
		return auth.NewDefaultUser("cs559", "1", nil, nil), nil
	}
	return nil, fmt.Errorf("invalid credentials")
}

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

func SetupGoGuardian() {
	authenticator = auth.New()
	cache = store.NewFIFO(context.Background(), time.Minute*5)
	basicStrategy := basic.New(ValidateUser, cache)
	tokenStrategy := bearer.New(VerifyToken, cache)
	authenticator.EnableStrategy(basic.StrategyKey, basicStrategy)
	authenticator.EnableStrategy(bearer.CachedStrategyKey, tokenStrategy)
}
