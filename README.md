# go-load-test
 Building a tool that can make a high volume of concurrent HTTP(S) request against an website/HTTP based API

## Build Command
```
go build cmd/bin/main.go
```

## Run Command
```
./main <arguments>
```
- #### Supported Arguments
  1. ```
      -u string
        URL to test
     ```
  2. ```
     -n int
        Number of requests to make (default 1)     
     ```
  3. ```
     -c int
        Number of concurrent requests to make (default 1)
     ```
  4. ```
     -f string
        File containing URLs to test (Currently on GET request is supported)
     ```
  5. ```
     -m string
        HTTP method for request (default "GET")
     ```
  6. ```
     -content-type string
        content-type for request (default "text/plain")
     ```
  7. ```
     -body string
        Body for the request
     ```
  8. ```
      -auth <bool>
        Authorization token is there ?
     ```
  9. ```
      -token string
        Authorization token
     ```
  10. ```
      -basic-auth <bool>
       Basic authentication has to be done ?
      ```
  11. ```
      -username string
        username for basic authentication
      ```
  12. ```
      -password string
        password for basic authentication
      ```

- For more details, visit [here](https://codingchallenges.fyi/challenges/challenge-load-tester/)
     
     
     
