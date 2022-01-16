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

EC2 instance link
```
curl http://ec2-65-0-130-180.ap-south-1.compute.amazonaws.com:8080/ --include --header "Content-Type: application/json" --request "POST" --data '{"long_url": "https://www.google.com", "short_url": ""}'
```

POST Request (to get custom short link):
```
curl http://localhost:8080/ --include --header "Content-Type: application/json" --request "POST" --data '{"longUrl": "https://www.google.com", "shortUrl": "custom-text"}'
```