package internals

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

type Result struct {
	statusCode      int
	totalTime       int
	timeToFirstByte int
}

type ExecuterTask struct {
	requestInfo *Request
}

func (et ExecuterTask) Execute(workerId int) TaskResult {
	startTime := time.Now()
	request := et.requestInfo

	body := bytes.NewReader(request.Body)
	req, err := http.NewRequest(request.Method, request.Url, body)

	if err != nil {
		log.Fatalf("Error : %s \n", err.Error())
	}

	if request.Auth.AuthType == BASIC {
		req.SetBasicAuth(request.Auth.Username, request.Auth.Password)
	} else if request.Auth.AuthType == TOKEN {
		req.Header.Set("Authorization", request.Auth.Token)
	}

	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("Error : %s \n", err.Error())
	}

	defer res.Body.Close()

	recievingFirst := time.Now()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatalf("Error : %s \n", err.Error())
	}

	log.Println(string(data))

	return Result{
		statusCode:      res.StatusCode,
		totalTime:       int(time.Since(startTime).Milliseconds()),
		timeToFirstByte: int(time.Since(recievingFirst).Milliseconds()),
	}
}
