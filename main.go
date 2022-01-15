package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

const hostUrl = "http://localhost:8080/"
const hostPort = "localhost:8080"

func sha256Of(input string) []byte {
	var algo = sha256.New()
	algo.Write([]byte(input))
	// log.Printf("%x\n", algo.Sum(nil))
	return algo.Sum(nil)
}

var repo = NewDynamoDBRepository()
var dynamoDBClient = createDynamoDBClient()

func Encode(msg string) string {
	urlHashBytes := sha256Of(msg)
	// println(urlHashBytes)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	// println(generatedNumber)
	encoded := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", generatedNumber)))
	return encoded[:6]
}

func saveUrlToDbandRespond(c *gin.Context, newUrlStruct Urlpair, shortUrl string) {
	var temp Urlpair = newUrlStruct
	temp.ShortURL = shortUrl
	newUrlStruct.ShortURL = hostUrl+shortUrl
	// TO DO: X = Check if longUrl already exists and return that if it does
	_,error := repo.Save(&temp)
	if error!=nil {
		// panic(error)		// TO DO remove panic and try to send error as string in json below
		log.Println(error)
		c.AbortWithStatus(500)
	} else {
		c.IndentedJSON(http.StatusCreated, newUrlStruct)
	}
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
		var shortUrl string = newUrlStruct.ShortURL
		
		pair,error := repo.FindByID(shortUrl)
		if error!=nil {
			log.Println(error)
			c.AbortWithStatus(500)
		} else {
			LongUrl := pair.LongURL
			switch LongUrl{
				case "" :
					saveUrlToDbandRespond(c,newUrlStruct, shortUrl)
				case newUrlStruct.LongURL:
					newUrlStruct.ShortURL = hostUrl+shortUrl
					c.IndentedJSON(http.StatusCreated, newUrlStruct)
				default:
					c.AbortWithStatusJSON(http.StatusConflict, gin.H{ "error" : "This Custom URL is not available"})
			}
		}

	} else {
		var shortUrl = Encode(newUrlStruct.LongURL)
		saveUrlToDbandRespond(c, newUrlStruct, shortUrl)
		// newUrlStruct.ShortURL = hostUrl+shortUrl
	}
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	pair,error := repo.FindByID(shortUrl)
	if error!=nil {
		log.Println(error)
		c.AbortWithStatus(500)	// TO DO error is currently returning empty, fix it
	} else {
		initialUrl := pair.LongURL
		if(initialUrl != ""){
			c.Redirect(302, initialUrl)
		} else {
			c.AbortWithStatus(404)
		}
	}
}

func main() {
	router := gin.Default()
	router.POST("/", PostUrl)
	router.GET("/:shortUrl", HandleShortUrlRedirect)
	router.Run(hostPort)
}
