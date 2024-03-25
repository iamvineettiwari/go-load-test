package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/iamvineettiwari/go-load-test/internals"
)

func getClientsFromFile(filename, contentType, username, password, token string, basicAuth, tokenAuth bool, count, concurrent int, clients []*internals.Client) ([]*internals.Client, error) {
	file, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	data := strings.Split(string(file), "\n")

	if len(data) == 0 {
		return nil, errors.New("no url found")
	}

	nRequest := max(1, count/len(data))

	for _, url := range data {
		headers := make(map[string]string)
		headers["Content-Type"] = contentType

		reqAuth, err := internals.NewRequestAuth(basicAuth, tokenAuth, username, password, token)

		if err != nil {
			return nil, err
		}

		clients = append(clients, internals.NewClient(
			nRequest,
			concurrent,
			internals.NewRequest(url, internals.GET, headers, []byte(nil), reqAuth),
		))
	}

	return clients, nil
}

func parseAndGetClient() ([]*internals.Client, error) {
	clients := []*internals.Client{}

	file := flag.String("f", "", "File containing URLs to test")
	url := flag.String("u", "", "URL to test")
	method := flag.String("m", "GET", "HTTP method for request")
	contentType := flag.String("content-type", "text/plain", "content-type for request")

	body := flag.String("body", "", "Body for the request")

	count := flag.Int("n", 1, "Number of requests to make")
	concurrent := flag.Int("c", 1, "Number of concurrent requests to make")

	basicAuth := flag.Bool("basic-auth", false, "Basic authentication has to be done ?")
	username := flag.String("username", "", "username for basic authentication")
	password := flag.String("password", "", "password for basic authentication")

	tokenAuth := flag.Bool("auth", false, "Authorization token is there ?")
	token := flag.String("token", "", "Authorization token")

	flag.Parse()

	if *count < 1 {
		return nil, errors.New("number of requests can not be less than 1")
	}

	if *concurrent < 1 {
		return nil, errors.New("concurrency can not be less than 1")
	}

	if *file != "" {
		return getClientsFromFile(*file, *contentType, *username, *password, *token, *basicAuth, *tokenAuth, *count, *concurrent, clients)
	}

	if *url == "" {
		return nil, errors.New("url is required")
	}

	headers := make(map[string]string)
	headers["Content-Type"] = *contentType

	reqAuth, err := internals.NewRequestAuth(*basicAuth, *tokenAuth, *username, *password, *token)

	if err != nil {
		return nil, err
	}

	clients = append(clients, internals.NewClient(
		*count,
		*concurrent,
		internals.NewRequest(*url, *method, headers, []byte(*body), reqAuth),
	))

	return clients, nil
}

func main() {
	clients, err := parseAndGetClient()

	if err != nil {
		log.Fatal(err.Error())
	}

	workerPool := internals.NewWorkerPool(max(1, min(4, len(clients))), len(clients))

	for _, client := range clients {
		workerPool.AddTask(client)
	}

	workerPool.Start()
	workerPool.Wait()
}
