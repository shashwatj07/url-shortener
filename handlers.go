package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Utility function to handle the logic of saving short links
// to the linked DynamoDB Instance along with the expiry time.
func saveUrlToDbandRespond(c *gin.Context, newUrlStruct urlStruct, shortUrl string) {
	var temp urlStruct = newUrlStruct
	temp.ShortURL = shortUrl
	newUrlStruct.ShortURL = HOST_URL + shortUrl
	days := newUrlStruct.ExpDate
	// Expiry date period has to be at least one day
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
		newUrlStruct.ExpDate = DefaultExpDate
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
	//check if ExpDate provided or not, if not set default
	if newUrlStruct.ExpDate == 0 {
		newUrlStruct.ExpDate = DefaultExpDate
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
			case newUrlStruct.LongURL:
				// If custom url is already allocated for same long url then return the same
				newUrlStruct.ShortURL = HOST_URL + shortUrl
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

// Handler to handle short URL's redirection.
func Redirect(c *gin.Context) {
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

// Handler to delete shortened urls from database before their expiry date
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
