package main

import (
	"os"
	"strconv"
)

// default target sec is not 0 to avoid the effect from cron/batch execution.
const defaultTargetSec = 31

func GetEnvTargetSec() int {
	num, err := strconv.Atoi(os.Getenv("METRIC_TARGET_SECOND"))
	if err != nil {
		return defaultTargetSec
	}
	switch {
	case num < 0,
		num > 59:
		return defaultTargetSec
	}
	return num
}

func GetEnvEventGroup() string {
	return os.Getenv("METRIC_EVENT_GROUP")
}

func GetEnvEventName() string {
	return os.Getenv("METRIC_EVENT_NAME")
}
