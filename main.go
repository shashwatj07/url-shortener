package main

import (
	"github.com/gin-gonic/gin"
)

// Entrypoint of the program
func main() {
	SetupGoGuardian()
	router := gin.Default()
	router.POST("/", AuthorizationMiddleware(), PostUrl)
	router.GET("/:shortUrl", Redirect)
	router.GET("/auth/token", AuthorizationMiddleware(), CreateToken)
	router.Run(HOST_PORT)
}
