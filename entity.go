package main

// ShortURL: The custom url or the hash we assign
// LongURL: The url to be shortened
// ExpDate: In the request it is the number of days for which the url is valid
// 		    In the response it is the UNIX epoch timestamp of the expiry date
type urlStruct struct {
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
	ExpDate  int    `json:"exp_date"`
}
