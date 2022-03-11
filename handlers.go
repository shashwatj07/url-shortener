package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var(
	addr  =  []string{"url-cache-0001-001.ncfh49.0001.aps1.cache.amazonaws.com:6379", "url-cache-0001-002.ncfh49.0001.aps1.cache.amazonaws.com:6379","url-cache-0001-003.ncfh49.0001.aps1.cache.amazonaws.com:6379"}
	postCache      PostCache           = NewRedisCache(addr, 0, 3600)
)

// Utility function to handle the logic of saving short links
// to the linked DynamoDB Instance along with the TTL.
func saveUrlToDbandRespond(c *gin.Context, newUrlStruct urlStruct, shortUrl string) {
	var temp urlStruct = newUrlStruct
	temp.ShortURL = shortUrl
	newUrlStruct.ShortURL = HOST_URL + shortUrl
	days := newUrlStruct.ExpDate
	// Validity period has to be at least one day
	if days < 1 {
		c.AbortWithStatus(500)
	} else {
		tempTime := time.Now().AddDate(0, 0, days)
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

func MustBindWith(c *gin.Context, newUrlStruct *urlStruct) error {

	if err := c.BindJSON(newUrlStruct); err != nil {
		return err
	}
	if newUrlStruct.ExpDate == 0 {
		newUrlStruct.ExpDate = DEFAULT_VALIDITY_DAYS
	}
	if newUrlStruct.LongURL == "" {
		return errors.New("long_url not provided or is empty")
	}
	return nil
}

// Handles the POST request to shorten a link. Performs the
// necessary sanity checks and properties to be followed according
// to defined conventions. Responds with the appropriate error
// or the short link in case of success.
func PostUrl(c *gin.Context) {
	var newUrlStruct urlStruct

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := MustBindWith(c, &newUrlStruct); err != nil {
		c.AbortWithStatusJSON(http.StatusExpectationFailed,
			gin.H{"error": err.Error()})
		return
	}
	// Check if ExpDate provided or not, if not set default
	if newUrlStruct.ExpDate == 0 {
		newUrlStruct.ExpDate = DEFAULT_VALIDITY_DAYS
	}
	// Add the new album to the slice.
	if newUrlStruct.ShortURL != "" {
		// custom shortURL
		if !IsAcceptableAlias(newUrlStruct.ShortURL) {
			c.AbortWithStatusJSON(http.StatusExpectationFailed,
				gin.H{"error": "Bad custom URL. Custom URL may only contain upto 32 letters, digits, underscore and hyphen symbols."})
		}

		var shortUrl string = newUrlStruct.ShortURL

		// Check if the given custom URL is available
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
				go saveUrltoAnalyticsDB(newUrlStruct, shortUrl)
			default:
				// If custom url is allocated to a long url return 409 Conflict status
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "This Custom URL is not available"})
			}
		}

	} else {
		//If no custom url provided then create a random short url using sha256 and then save to database
		var shortUrl = Encode(newUrlStruct.LongURL)
		saveUrlToDbandRespond(c, newUrlStruct, shortUrl)
		go saveUrltoAnalyticsDB(newUrlStruct, shortUrl)
	}
}

// Handler to handle short URL's redirection.
func Redirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	longURL, err := postCache.Get(shortUrl)
	if err !=nil{
		pair, error := repo.FindByID(shortUrl)
		if error != nil {
			log.Println(error)
			c.AbortWithStatus(500)
		} else {
			longURL = pair.LongURL
			if longURL != "" {
				// Redirect to original url
				c.Redirect(302, longURL)
				postCache.Set(shortUrl,longURL)
				go incrementRedirCount(shortUrl)
			} else {
				// Short url does not exist
				c.AbortWithStatus(404)
			}
		}
	} else {
		c.Redirect(302, longURL)
		go incrementRedirCount(longURL)
	}
	
}

// Get Analytics for a url based on per day usage
//
// Returns a JSON response containing date and usage on that date if url is found else 404
func GetAnalytics(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	analytics, error := GetAnalyticsFromDb(shortUrl)
	if error != nil {
		log.Println(error)
		c.AbortWithStatus(500)
	} else {
		if analytics != nil {
			// Short url exists then return the analytics
			c.IndentedJSON(http.StatusFound, analytics)
		} else {
			// Short url does not exist
			c.AbortWithStatus(404)
		}
	}
}

// Handler to delete shortened urls from database before their validity ends
func DeleteUrl(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	error := repo.Delete(shortUrl)
	if error != nil {
		log.Println(error)
		c.AbortWithStatus(500)
	} else {
		c.Status(204)
	}
}

// Middleware function to intercept the API request and
// check if it is authorized to proceed. Aborts the request
// if it is found to be unauthorized.
func AuthorizationMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		var r *http.Request = c.Request
		log.Println("Executing Auth Middleware")
		user, err := authenticator.Authenticate(r)
		if err != nil {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			c.Abort()
			return
		}
		log.Printf("User %s Authenticated\n", user.UserName())
		c.Next()
	})
}

// Handler for bulk URL shortening from CSV.
func PostBulkUrl(c *gin.Context) {
	header, receiveErr := c.FormFile("file")
	if receiveErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": "CSV file not provided."})
	}
	fmt.Println(header.Filename)
	out, openErr := header.Open()
	if openErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": "Unable to open file."})
	}
	csvLines, readErr := csv.NewReader(out).ReadAll()
	if readErr != nil {
		c.AbortWithStatusJSON(http.StatusExpectationFailed,
			gin.H{"error": "File must be a CSV."})
	}
	responses := make([]urlStruct, len(csvLines))
	for index, line := range csvLines {
		longUrl := line[0]
		alias := line[1]
		validity, err := strconv.Atoi(line[2])
		if err != nil {
			validity = 30
		}
		responses[index] = PostUrlUtil(longUrl, alias, validity)
		go saveUrltoAnalyticsDB(responses[index], responses[index].ShortURL[strings.LastIndex(responses[index].ShortURL, "/")+1:])
	}
	c.IndentedJSON(http.StatusAccepted, responses)
}
