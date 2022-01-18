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
(Note: The expiry period is 30 days for the bearer token generated)

GET Request (to load shortened url):
```
curl http://localhost:8080/<random-hash|custom-alias>
```

POST Request (to get random short link):
```
curl http://localhost:8080/ --include --header "Content-Type: application/json" --request "POST" --data '{"long_url": "https://www.google.com", "short_url": "", "exp_date": 31}' -H "Authorization: Bearer <auth_token>"
```

POST Request (to get custom short link):
```
curl http://localhost:8080/ --include --header "Content-Type: application/json" --request "POST" --data '{"long_url": "https://www.google.com", "short_url": "custom-text", "exp_date": 31}' -H "Authorization: Bearer <auth_token>"
```
(Note: The default expiry period is 30 days for the shortened URL if not specified in the request)

DELETE Request (to delete a short link before it expires)
```
curl http://localhost:8080/<random-hash|custom-alias> --request DELETE -H "Authorization: Bearer <auth_token>"
```

(Note: Replace localhost with public DNS entry for accessing hosted version on AWS. Current public DNS: ec2-65-0-130-180.ap-south-1.compute.amazonaws.com)
