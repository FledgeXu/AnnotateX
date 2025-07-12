package utils

import (
	"strconv"
	"time"
)

func UnixStringToTime(s string) (time.Time, error) {
	unixSec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(unixSec, 0), nil
}
