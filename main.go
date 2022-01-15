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

func sha256Of(input string) []byte {
	var algo = sha256.New()
	algo.Write([]byte(input))
	return algo.Sum(nil)
}

type urlStruct struct {
	LongURL  string `json:"longUrl"`
	ShortURL string `json:"shortUrl"`

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
		var temp Urlpair = newUrlStruct
		temp.ShortURL = shortUrl
		newUrlStruct.ShortURL = hostUrl+shortUrl
		// TO DO: X = Check if longUrl already exists and return that if it does
		_,error := repo.Save(&temp)
		if error!=nil {
			panic(error)		// TO DO remove panic and try to send error as string in json below
			c.AbortWithStatusJSON(500, gin.H{"error": error})
		} else {
			c.IndentedJSON(http.StatusCreated, newUrlStruct)
		}
		newUrlStruct.ShortURL = hostUrl+shortUrl
	}
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	pair,error := repo.FindByID(shortUrl)
	if error!=nil {
		c.AbortWithStatusJSON(500, gin.H{"error": error})	// TO DO error is currently returning empty, fix it
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
