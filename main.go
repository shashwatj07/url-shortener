package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

const hostUrl = "http://localhost:8080/"
const hostPort = "localhost:8080"

func sha256Of(input string) []byte {
	/*
		@param input : url to be shortened
		return: shasum of the sha256 hash
	*/
	var algo = sha256.New()
	algo.Write([]byte(input))
	return algo.Sum(nil)
}

var repo = NewDynamoDBRepository()
var dynamoDBClient = createDynamoDBClient()

func Encode(msg string) string {
	/*
		@param msg : the url that needs to be shortened
		Create sha256 hash of the long url use its first six characters as the short url
		return: first six characters as the short url
	*/
	urlHashBytes := sha256Of(msg)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	encoded := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", generatedNumber)))
	return encoded[:6]
}

func saveUrlToDbandRespond(c *gin.Context, newUrlStruct Urlpair, shortUrl string) {
	/*
		@param c : the context propagated from the router
		@param newUrlStruct : (shortUrl, longUrl, ExpDate) To be returned in the response
		@param shortUrl : Carries the custom or random 6 character hash
		Checks if the exp_date period provided by user is at least one day. If not 
		it returns bad request else it creates a new short-long url pair in the database table
		return: Nothing
	*/
	var temp Urlpair = newUrlStruct
	temp.ShortURL = shortUrl
	newUrlStruct.ShortURL = hostUrl + shortUrl
	days := newUrlStruct.ExpDate
	// Expiry date period has to be at least one day
	if days<1 {
		c.AbortWithStatus(500)
	} else {
		tempTime := time.Now().AddDate(0,0,days)
		temp.ExpDate = int(tempTime.Unix()) // Get the UNIX epoch timestamp
		_, error := repo.Save(&temp)
		if error != nil {
			log.Println(error)
			c.AbortWithStatus(500)
		} else {
			// New entry successful
			c.IndentedJSON(http.StatusCreated, newUrlStruct)
		}
	}
}

func PostUrl(c *gin.Context) {
	/*
		@param c : the context propagated from the router
		Check if user has provided a custom url. if the custom url is available then provide that as the url
		else return 409 Conflict status to indicate url is not available. If custom url not provided then a random one is created using the sha256 hash function
		return: Nothing
	*/
	var newUrlStruct Urlpair

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newUrlStruct); err != nil {
		return
	}
	// Add the new album to the slice.
	if newUrlStruct.ShortURL != "" {
		//custom shortURL
		var shortUrl string = newUrlStruct.ShortURL

		//Check if the given custom URL is available
		pair, error := repo.FindByID(shortUrl)
		if error != nil {
			log.Println(error)
			c.AbortWithStatus(500)
		} else {
			LongUrl := pair.LongURL
			switch LongUrl {
			case "":
				// If custom url is available create a new entry with it
				saveUrlToDbandRespond(c, newUrlStruct, shortUrl)
			case newUrlStruct.LongURL:
				// If custom url is already allocated for same long url then return the same
				newUrlStruct.ShortURL = hostUrl + shortUrl
				c.IndentedJSON(http.StatusCreated, newUrlStruct)
			default:
				// If custom url is allocated to different long url return 409 Conflict status
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "This Custom URL is not available"})
			}
		}

	} else {
		//If no custom url provided then create a random short url using sha256 and then save to database
		var shortUrl = Encode(newUrlStruct.LongURL)
		saveUrlToDbandRespond(c, newUrlStruct, shortUrl)
	}
}

func HandleShortUrlRedirect(c *gin.Context) {
	/*
		@param c : the context propagated from the router
		Check if there is a long url for corresponding short url. If it exists then redirect to it else 404
		return: Nothing
	*/
	shortUrl := c.Param("shortUrl")
	pair, error := repo.FindByID(shortUrl)
	if error != nil {
		log.Println(error)
		c.AbortWithStatus(500) 
	} else {
		initialUrl := pair.LongURL
		if initialUrl != "" {
			//redirect to origignal url
			c.Redirect(302, initialUrl)  
		} else {
			// Short url does not exist
			c.AbortWithStatus(404) 
		}
	}
}

func main() {
	// Setup the router
	router := gin.Default()
	router.POST("/", PostUrl)
	router.GET("/:shortUrl", HandleShortUrlRedirect)
	router.Run(hostPort)
}
