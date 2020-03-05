package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/evalphobia/aws-sdk-go-wrapper/cloudtrail"
	"github.com/evalphobia/aws-sdk-go-wrapper/cloudwatch"
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
)

func getAndSendMetrics(events []string, timePoints []time.Time) error {
	groupName := GetEnvEventGroup()
	if groupName == "" {
		return errors.New("EventGroup is empty")
	}

	trailCli, err := cloudtrail.New(config.Config{})
	if err != nil {
		return err
	}
	cwCli, err := cloudwatch.New(config.Config{})
	if err != nil {
		return err
	}

	attrs := make([]cloudtrail.LookupAttribute, len(events))
	for i, ev := range events {
		attrs[i] = cloudtrail.LookupAttribute{
			Key:   "EventName",
			Value: ev,
		}
	}

	// get metrics
	metrics := make([]cloudwatch.MetricDatum, 0, len(events)*len(timePoints))
	for _, dt := range timePoints {
		resp, err := trailCli.LookupEventsAll(cloudtrail.LookupEventsInput{
			StartTime:        dt,
			EndTime:          dt,
			LookupAttributes: attrs,
		})
		if err != nil {
			return err
		}

		eventCounts := make(map[string]int)
		for _, ev := range resp.Events {
			eventCounts[ev.EventName]++
		}

		for _, ev := range events {
			count, _ := eventCounts[ev]
			metrics = append(metrics, cloudwatch.MetricDatum{
				MetricName: fmt.Sprintf("%s::%s", groupName, ev),
				Value:      float64(count),
				HasValue:   true,
				Timestamp:  dt,
				Unit:       "Count/Second",
			})
		}
	}

	// send metrics
	in := cloudwatch.PutMetricDataInput{
		Namespace:  "AWSAPI/" + groupName,
		MetricData: metrics,
	}

	// logging
	j, _ := json.Marshal(in.ToInput())
	fmt.Printf("---- [PutMetricDataInput] ----\n%s\n---- ---- ---- ----\n", string(j))

	return cwCli.PutMetricData(in)
}
