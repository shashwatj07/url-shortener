package main

// import "encoding/json"

// "time"

type Urlpair struct {
	ShortURL    string  `json:"short_url"`
	LongURL  string `json:"long_url"`
	ExpDate int `json:"exp_date"`
}
