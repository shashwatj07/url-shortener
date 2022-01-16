package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/url"
	"regexp"
)

var HasValidCustomLinkChars = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString

func sha256Of(input string) []byte {
	var algo = sha256.New()
	algo.Write([]byte(input))
	return algo.Sum(nil)
}

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

func IsAcceptableAlias(alias string) bool {
	return HasValidCustomLinkChars(alias) && len(alias) <= 32
}
