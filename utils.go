package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/url"
	"regexp"
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
