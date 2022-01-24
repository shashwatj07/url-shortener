package main

// ShortURL: The custom url or the hash we assign
//
// LongURL: The url to be shortened
//
// ExpDate: In the request it is the number of days after which the url expires
//
// 		    In the response it is the UNIX epoch timestamp of the expiry date
type urlStruct struct {
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
	ExpDate  int    `json:"exp_date"`
}

// Date: The date when a url was visited
//
// Count: The number of times it was visited that day
type DateCountStruct struct {
    Date string	`json:"date"`
    Count int	`json:"count"`
}
