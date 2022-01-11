# url-shortener

Fire up the server:
```
go get .
go run main.go
```

POST Request (to get random short link):
```
curl http://localhost:8080/ --include --header "Content-Type: application/json" --request "POST" --data '{"longUrl": "https://www.google.com", "shortUrl": ""}'
```

POST Request (to get custom short link):
```
curl http://localhost:8080/ --include --header "Content-Type: application/json" --request "POST" --data '{"longUrl": "https://www.google.com", "shortUrl": "custom-text"}'
```
