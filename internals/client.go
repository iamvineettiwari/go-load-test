package internals

import (
	"fmt"
	"math"
	"time"
)

type ResultStatistics struct {
	url              string
	totalTimeForTest float32
	totalRequests    int
	totalSuccess     int
	totalFailed      int
	requestPerSecond float32
	requestTimeMin   float32
	requestTimeMax   float32
	requestTimeAvg   float32
	firstByteTimeMin float32
	firstByteTimeMax float32
	firstByteTimeAvg float32
}

type Client struct {
	workerPool   *Pool
	totalRequest int
	concurrency  int
	requestInfo  *Request
}

func NewClient(totalRequest, concurrency int, requestInfo *Request) *Client {
	return &Client{
		totalRequest: totalRequest,
		concurrency:  concurrency,
		requestInfo:  requestInfo,
		workerPool:   NewWorkerPool(concurrency, totalRequest),
	}
}

func (c *Client) Execute(clientId int) TaskResult {
	c.Test()
	return nil
}

func (c *Client) Test() {
	testStartTime := time.Now()

	for taskId := 1; taskId <= c.totalRequest; taskId++ {
		c.workerPool.AddTask(ExecuterTask{
			requestInfo: c.requestInfo,
		})
	}

	c.workerPool.Start()
	c.workerPool.Wait()
	results := c.workerPool.CollectResult()

	c.calculateAndPrintStatistics(testStartTime, results)
}

func (c *Client) calculateAndPrintStatistics(testStartTime time.Time, results []TaskResult) {
	totalSuccess, totalError := 0, 0
	requestTimeMin, requestTimeMax, totalTime := math.MaxInt, math.MinInt, 0
	timeToFirstByteMin, timeToFirstByteMax, totalByteTime := math.MaxInt, math.MinInt, 0

	for _, resultItem := range results {
		result := resultItem.(Result)

		if result.statusCode >= 200 && result.statusCode <= 299 {
			totalSuccess++
		}

		if result.statusCode >= 400 && result.statusCode <= 599 {
			totalError++
		}

		requestTimeMin = min(requestTimeMin, result.totalTime)
		requestTimeMax = max(requestTimeMax, result.totalTime)
		totalTime += result.totalTime

		timeToFirstByteMin = min(timeToFirstByteMin, result.timeToFirstByte)
		timeToFirstByteMax = max(timeToFirstByteMax, result.timeToFirstByte)
		totalByteTime += result.timeToFirstByte

	}

	totalTimeForTest := float32(time.Since(testStartTime).Milliseconds()) / 1000
	avgRequestTime := (float32(totalTime) / 1000) / float32(c.totalRequest)
	avgFirstByteTime := (float32(totalByteTime) / 1000) / float32(c.totalRequest)
	requestPerSecond := float32(c.totalRequest) / (float32(totalTime) / 1000)

	rs := ResultStatistics{
		url:              c.requestInfo.Url,
		totalRequests:    c.totalRequest,
		totalTimeForTest: totalTimeForTest,
		totalSuccess:     totalSuccess,
		totalFailed:      totalError,
		requestPerSecond: requestPerSecond,
		requestTimeMin:   float32(requestTimeMin) / 1000,
		requestTimeMax:   float32(requestTimeMax) / 1000,
		requestTimeAvg:   avgRequestTime,
		firstByteTimeMin: float32(timeToFirstByteMin) / 1000,
		firstByteTimeMax: float32(timeToFirstByteMax) / 1000,
		firstByteTimeAvg: avgFirstByteTime,
	}

	rs.Print()
}

func (rs ResultStatistics) Print() {
	fmt.Printf("\nResults for : %s\n", rs.url)
	fmt.Println()
	fmt.Printf("\t Total Time Taken for test (s)         : %f\n", rs.totalTimeForTest)
	fmt.Printf("\t Total Requests                        : %d\n", rs.totalRequests)
	fmt.Printf("\t Total Succeed (2XX)                   : %d\n", rs.totalSuccess)
	fmt.Printf("\t Total Failed (4XX-5XX)                : %d\n", rs.totalFailed)
	fmt.Printf("\t Resquest per second                   : %f\n", rs.requestPerSecond)
	fmt.Println()
	fmt.Printf("\tRequest Time (s) (min, max, avg)       : %f, %f, %f\n", rs.requestTimeMin, rs.requestTimeMax, rs.requestTimeAvg)
	fmt.Printf("\tTime to First Byte (s) (min, max, avg) : %f, %f, %f\n", rs.firstByteTimeMin, rs.firstByteTimeMax, rs.firstByteTimeAvg)
	fmt.Println()
}
