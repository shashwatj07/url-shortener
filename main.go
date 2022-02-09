package main

import (
	"github.com/gin-gonic/gin"
)

// Entrypoint of the program
func main() {
	SetupGoGuardian()
	SetupSqlDbConnection()
	router := gin.Default()
	router.POST("/", AuthorizationMiddleware(), PostUrl)
	router.DELETE("/:shortUrl", AuthorizationMiddleware(), DeleteUrl)
	router.GET("/:shortUrl", Redirect)
	router.GET("/analytics/:shortUrl", AuthorizationMiddleware(), GetAnalytics)
	router.GET("/auth/token", AuthorizationMiddleware(), CreateToken)
	router.POST("/bulk", AuthorizationMiddleware(), PostBulkUrl)
	router.Run(HOST_PORT)
}
