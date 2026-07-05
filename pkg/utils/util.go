package utils

import (
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func NewError(message string, args ...any) error {
	if args != nil {
		return fmt.Errorf("%s", fmt.Sprintf(message, args...))
	}
	return fmt.Errorf("%s", message)
}

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

func StrToUint(val string) uint {
	var uintValue uint
	_, err := fmt.Sscanf(val, "%d", &uintValue)
	if err != nil {
		return 0
	}
	return uintValue
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

// Struct to map
func StructToMap(data any) (map[string]any, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
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

func NormalizeEntityName(name string) string {
	if name == "" {
		return name
	}
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "_", " ")
	return strings.ReplaceAll(name, "-", " ")
}

func TransformPositionLabel(position *string) {
	if position == nil || *position == "" {
		return
	}
	*position = Position(*position).Label()
}

func IsValidEmail(email string) bool {
	// Simple regex for email validation
	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

func IsValidMobile(mobile string) bool {
	// Simple regex for mobile number validation (10 digits)
	const mobileRegex = `^\d{10}$`
	return regexp.MustCompile(mobileRegex).MatchString(mobile)
}

func GenrateOTP(n int) string {
	if n <= 0 {
		n = 4
	}

	max := int64(1)
	for i := 0; i < n; i++ {
		max *= 10
	}

	num, _ := crand.Int(crand.Reader, big.NewInt(max))
	return fmt.Sprintf("%0*d", n, num.Int64())
}

func HashOTP(otp, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(otp))
	return hex.EncodeToString(h.Sum(nil))
}

func SecureCompare(a, b string) bool {
	return hmac.Equal([]byte(a), []byte(b))
}

func ClientIP(c *gin.Context) string {
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	ip := c.ClientIP()
	p := net.ParseIP(ip)
	if p == nil {
		return ip
	}
	// IPv4-mapped IPv6 (::ffff:1.2.3.4) -> return IPv4
	if p4 := p.To4(); p4 != nil {
		return p4.String()
	}
	// map IPv6 loopback to IPv4 loopback
	if p.IsLoopback() {
		return "127.0.0.1"
	}
	return p.String()
}

func Contains[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}
