package main

type PostCache interface {
	Set(key string, value string)
	Get(key string) (string,error)
}
