package utils

import (
	"encoding/json"
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

func MapToStruct(m map[string]interface{}, result interface{}) error {
	for k, v := range m {
		if str, ok := v.(string); ok {
			var tmp interface{}
			if err := json.Unmarshal([]byte(str), &tmp); err == nil {
				m[k] = tmp
			}
		} else if bytes, ok := v.([]byte); ok {
			var tmp interface{}
			if err := json.Unmarshal(bytes, &tmp); err == nil {
				m[k] = tmp
			}
		}
	}
	// Convert map to JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Convert JSON to struct
	return json.Unmarshal(jsonData, result)
}

func StructToStr(data any) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func StrToBool(s string) bool {
	return s == "true" || s == "1"
}
