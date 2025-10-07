package utils

import (
	"fmt"
	"time"
)

func CreateCurrentDateTime() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

func StrToInt(val string) (int, error) {
	var intValue int
	_, err := fmt.Sscanf(val, "%d", &intValue)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}
