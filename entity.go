package main

type urlStruct struct {
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
	ExpDate  int    `json:"exp_date"`
}
