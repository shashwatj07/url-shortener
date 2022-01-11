package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"github.com/gin-gonic/gin"
)

const host = "http://localhost:8080/"

var algo = sha256.New()

func sha256Of(input string) []byte {
	algo.Write([]byte(input))
	return algo.Sum(nil)
}

type urlStruct struct {
	LongURL  string `json:"longUrl"`
	ShortURL string `json:"shortUrl"`
}

var store = make(map[string]string)

func Encode(msg string) string {
	urlHashBytes := sha256Of(msg)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	encoded := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", generatedNumber)))
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
		newUrlStruct.ShortURL = host+shortUrl
		store[shortUrl] = newUrlStruct.LongURL
	}

	c.IndentedJSON(http.StatusCreated, newUrlStruct)
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	initialUrl := store[shortUrl]
	c.Redirect(302, initialUrl)
}

func main() {
	router := gin.Default()
	router.POST("/", PostUrl)
	router.GET("/:shortUrl", HandleShortUrlRedirect)
	router.Run("localhost:8080")
}
