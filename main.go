package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event LambdaEvent) (LambdaResponse, error) {
	events := event.getEventNameList()
	dateList := event.getTimePointList()
	switch {
	case len(events) == 0,
		len(dateList) == 0:
		err := errors.New("Target values are empty")
		return NewErrResponse(err), err
	}

	// fetch data from CloudTrail events and send them to CloudWatch Metrics.
	err := getAndSendMetrics(events, dateList)
	if err != nil {
		return NewErrResponse(err), err
	}

	return LambdaResponse{
		Status:  200,
		Success: true,
	}, nil
}

type LambdaResponse struct {
	Status  int   `json:"status"`
	Success bool  `json:"success"`
	Err     error `json:"error"`
}

func NewErrResponse(err error) LambdaResponse {
	return LambdaResponse{
		Status:  500,
		Success: false,
		Err:     err,
	}
}

type LambdaEvent struct {
	EventName string `json:"event_name"`
}

func (e LambdaEvent) getEventNameList() []string {
	n := e.EventName
	if n == "" {
		n = GetEnvEventName()
	}
	list := strings.Split(n, ",")
	for i, v := range list {
		list[i] = strings.TrimSpace(v)
	}
	return list
}

func (e LambdaEvent) getTimePointList() []time.Time {
	const baseDelayMin = 16 // CloudTrail has 15min delay in maximum
	const timePointsCount = 5

	now := time.Now()
	base := modifySec(now)

	list := make([]time.Time, timePointsCount)
	for i := range list {
		pastMin := -1 * time.Duration(baseDelayMin+i)
		list[i] = base.Add(time.Minute * pastMin)
	}

	return list
	// (e.g.)
	// return []time.Time{
	// 	base.Add(time.Minute * -16),
	// 	base.Add(time.Minute * -17),
	// 	base.Add(time.Minute * -18),
	// 	base.Add(time.Minute * -19),
	// 	base.Add(time.Minute * -20),
	// }
}

// change the second of given time.
// default is `31`.
func modifySec(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), GetEnvTargetSec(), 0, dt.Location())
}
