package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"regexp"
	"time"
)

// Regexp object to match valid characters in an alias.
var HasValidCustomLinkChars = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString

// Function to compute the SHA256 hash of the given string.
func sha256Of(input string) []byte {
	var algo = sha256.New()
	algo.Write([]byte(input))
	return algo.Sum(nil)
}

// Utility function to encode a long URL into a short random
// alias.
func Encode(msg string) string {
	urlHashBytes := sha256Of(msg)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	encoded := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", generatedNumber)))
	return encoded[:6]
}

// Utility to check is the given string is a valid URL.
func IsUrl(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Utility function to check whether the custom alias given is
// acceptable or not. By convention, it may only contain upto 32
// letters, digits, underscore and hyphen symbols.
func IsAcceptableAlias(alias string) bool {
	return HasValidCustomLinkChars(alias) && len(alias) <= 32
}

func SaveUrl(longUrl string, alias string, validity int) string {
	var temp urlStruct
	temp.ShortURL = alias
	temp.LongURL = longUrl
	// Expiry date period has to be at least one day
	if validity < 1 {
		return "ERROR: Validity cannot be less than 1 day."
	} else {
		tempTime := time.Now().AddDate(0, 0, validity)
		temp.ExpDate = int(tempTime.Unix()) // Get the UNIX epoch timestamp
		_, error := repo.Save(&temp)
		if error != nil {
			log.Println(error)
			return "ERROR: Failed to save url to DB."
		} else {
			// New entry successful
			return HOST_URL + alias
		}
	}
}

// Utility function to handle one item from a bulk shortening request.
// Returns a urlStruct object with necessary details.
func PostUrlUtil(longUrl string, alias string, validity int) (result urlStruct) {
	result.LongURL = longUrl
	result.ExpDate = validity
	if alias != "" {
		// custom shortURL
		if !IsAcceptableAlias(alias) {
			result.ShortURL = "ERROR: Bad custom URL. Custom URL may only contain upto 32 letters, digits, underscore and hyphen symbols."
		} else {
			// Check if the given custom URL is available
			pair, error := repo.FindByID(alias)
			if error != nil {
				log.Println(error)
				result.ShortURL = "ERROR: Query Failed."
			} else {
				switch pair.LongURL {
				case "":
					// If custom url is available create a new entry with it
					go SaveUrl(longUrl, alias, validity)
					result.ShortURL = HOST_URL + alias
				case longUrl:
					// If custom url is already allocated for same long url then return the same
					result.ShortURL = HOST_URL + alias
				default:
					// If custom url is allocated to different long url
					result.ShortURL = "ERROR: Requested Custom URL is not available"
				}
			}
		}
	} else {
		// If no custom url provided then create a random short url using sha256 and then save to database
		alias = Encode(longUrl)
		go SaveUrl(longUrl, alias, validity)
		result.ShortURL = HOST_URL + alias
	}
	return result
}
