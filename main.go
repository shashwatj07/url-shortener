package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const hostUrl = "http://localhost:8080/"
const hostPort = "localhost:8080"

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

func IsUrl(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func PostUrl(c *gin.Context) {
	var newUrlStruct urlStruct

	// Call BindJSON to bind the received JSON to newUrlStruct.
	if err := c.BindJSON(&newUrlStruct); err != nil {
		return
	}

	if !IsUrl(newUrlStruct.LongURL) {
		// Malformed URL
		newUrlStruct.ShortURL = "Malformed URL cannot be shortened."
	} else if newUrlStruct.ShortURL != "" {
		// Custom short link
		store[newUrlStruct.ShortURL] = newUrlStruct.LongURL
	} else {
		// Random short link
		var shortUrl = Encode(newUrlStruct.LongURL)
		newUrlStruct.ShortURL = hostUrl + shortUrl
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
	router.Run(hostPort)
}
