package webhook

import (
	"strconv"
	"time"
)

const WebhookTimestampTolerance = 5 * time.Minute

func VerifyTimestamp(timestamp string) bool {
	unixTimestamp, err := strconv.ParseInt(timestamp, 10, 64)

	if err != nil {
		return false
	}

	eventTime := time.Unix(unixTimestamp, 0)
	now := time.Now()

	difference := now.Sub(eventTime)

	if difference < 0 {
		difference = -difference
	}

	return difference <= WebhookTimestampTolerance
}
