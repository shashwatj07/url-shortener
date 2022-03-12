# url-shortener

Fire up the server-
```
go get .
go run .
```

GET Request (to generate Authorization Token):
```
curl http://url-shortener4-dev.ap-south-1.elasticbeanstalk.com/auth/token -u cs559:iitbh
```
(Note: The validity period is 30 days for the bearer token generated)

GET Request (to get analytics for a url):
```
curl http://localhost:8080/analytics/<random-hash|custom-alias> -H "Authorization: Bearer <auth_token>"
```

GET Request (to load shortened url):
```
curl http://localhost:8080/<random-hash|custom-alias>
```

POST Request (to get random short link):
```
curl http://url-shortener4-dev.ap-south-1.elasticbeanstalk.com/ --include --header "Content-Type: application/json" --request "POST" --data '{"long_url": "https://www.linkedin.com", "short_url": "", "exp_date": 31}' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2NDk2NTg0NjYsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoiY3M1NTkifQ._JshleXML9zqsV4sDAtaoBhxKPldyE2MPw_2Fo8XjGw"
```

POST Request (to get custom short link):
```
curl http://localhost:8080/ --include --header "Content-Type: application/json" --request "POST" --data '{"long_url": "https://www.google.com", "short_url": "custom-text", "exp_date": 31}' -H "Authorization: Bearer <auth_token>"
```
(Note: The default validity period is 30 days for the shortened URL if not specified in the request)

DELETE Request (to delete a short link before it expires)
```
curl http://url-shortener3-dev.ap-south-1.elasticbeanstalk.com/lkd --request DELETE -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2NDgxODgwMTgsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoiY3M1NTkifQ.WEYTfr9A0yndVIQRWBX4vgD7n6tkvQX2nh7eU1DgtLE"
```

POST Request (bulk shortening):
```
curl http://url-shortener3-dev.ap-south-1.elasticbeanstalk.com/bulk --request POST -F file="@test/test.csv" -H "Content-Type: multipart/form-data" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2NDgxODgwMTgsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoiY3M1NTkifQ.WEYTfr9A0yndVIQRWBX4vgD7n6tkvQX2nh7eU1DgtLE"
```

(Note: Replace localhost with public DNS entry for accessing hosted version on AWS. Current public DNS: ec2-65-0-130-180.ap-south-1.compute.amazonaws.com)
