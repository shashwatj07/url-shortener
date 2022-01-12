package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
)

const hostUrl = "http://localhost:8080/"
const hostPort = "localhost:8080"

var algo = sha256.New()

func sha256Of(input string) []byte {
	algo.Write([]byte(input))
	return algo.Sum(nil)
}

var repo = NewDynamoDBRepository()
var dynamoDBClient = createDynamoDBClient()

func Encode(msg string) string {
	urlHashBytes := sha256Of(msg)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	encoded := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", generatedNumber)))
	return encoded[:6]
}

func PostUrl(c *gin.Context) {
	var newUrlStruct Urlpair

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newUrlStruct); err != nil {
		return
	}
	// Add the new album to the slice.
	if newUrlStruct.ShortURL != "" {
		//TO DO
		//custom shortURL
	} else {
		var shortUrl = Encode(newUrlStruct.LongURL)
		newUrlStruct.ShortURL = hostUrl+shortUrl
		repo.Save(&newUrlStruct)
	}

	c.IndentedJSON(http.StatusCreated, newUrlStruct)
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	pair,error := repo.FindByID(shortUrl)
	if error!=nil {
		c.AbortWithStatusJSON(500, gin.H{"error": error})
	} else {
		initialUrl := pair.LongURL
		c.Redirect(302, initialUrl)
	}
}

func main() {
	router := gin.Default()
	router.POST("/", PostUrl)
	router.GET("/:shortUrl", HandleShortUrlRedirect)
	router.Run(hostPort)
}
