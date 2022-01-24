# url-shortener

Fire up the server:
```
go get .
go run .
```

GET Request (to generate Authorization Token):
```
curl http://localhost:8080/auth/token -u username:password
```

GET Request (to get anlytics for a url):
```
curl http://localhost:8080/analytics/<random-hash|custom-alias>
```

POST Request (to get random short link):
```
curl http://localhost:8080/ --include --header "Content-Type: application/json" --request "POST" --data '{"long_url": "https://www.google.com", "short_url": "", "exp_date": 30}' -H "Authorization: Bearer <auht_token>"
```

POST Request (to get custom short link):
```
curl http://localhost:8080/ --include --header "Content-Type: application/json" --request "POST" --data '{"long_url": "https://www.google.com", "short_url": "custom-text", "exp_date": 30}' -H "Authorization: Bearer <auth_token>"
```
