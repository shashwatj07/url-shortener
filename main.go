package main

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type urlStruct struct {
	LongURL  string `json:"longUrl"`
	ShortURL string `json:"shortUrl"`
}

var store = make(map[string]string)

func Encode(msg string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(msg))
	fmt.Println(encoded)
	return encoded[:6]
}

func PostUrl(c *gin.Context) {
	var newUrlStruct urlStruct

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newUrlStruct); err != nil {
		return
	}

	// Add the new album to the slice.
	if newUrlStruct.ShortURL != "" {
		store[newUrlStruct.ShortURL] = newUrlStruct.LongURL
	} else {
		var shortUrl = Encode(newUrlStruct.LongURL)
		newUrlStruct.ShortURL = shortUrl
		store[shortUrl] = newUrlStruct.LongURL
	}

	c.IndentedJSON(http.StatusCreated, newUrlStruct)
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	initialUrl := store[shortUrl]
	fmt.Println(initialUrl)
	c.Redirect(302, initialUrl)
}

func main() {
	router := gin.Default()
	router.POST("/", PostUrl)
	router.GET("/:shortUrl", HandleShortUrlRedirect)
	router.Run("localhost:8080")
}
