# url-shortener
Requirements:
```
AWSCLI
aws-sdk-go
```


Fire up the server:
```
go get .
go run entity.go repository.go main.go
```

POST Request (to get random short link):
```
curl http://localhost:8080/ --include --header "Content-Type: application/json" --request "POST" --data '{"longUrl": "https://www.google.com", "shortUrl": ""}'
```

POST Request (to get custom short link):
```
curl http://localhost:8080/ --include --header "Content-Type: application/json" --request "POST" --data '{"longUrl": "https://www.google.com", "shortUrl": "custom-text"}'
```

