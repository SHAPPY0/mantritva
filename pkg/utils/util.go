package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
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

func CompareIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	ac := append([]int(nil), a...)
	bc := append([]int(nil), b...)
	sort.Ints(ac)
	sort.Ints(bc)
	for i := range ac {
		if ac[i] != bc[i] {
			return false
		}
	}
	return true
}

func IsEmpty(s interface{}) bool {
	return reflect.DeepEqual(s, reflect.Zero(reflect.TypeOf(s)).Interface())
}

func RoundToOneDecimal(val float64) float64 {
	if val == 0 {
		return 0
	}
	rounded := math.Round(val*10) / 10
	if rounded == 0 {
		return 0
	}

	return rounded
}
