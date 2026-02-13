package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	// "github.com/shappy0/mantritva/pkg/config"
)

func CreateJWTToken(data map[string]interface{}) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	for k, v := range data {
		atClaims[k] = v
	}
	atClaims["exp"] = time.Now().Add(time.Minute * 60 * 24 * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	// token, err := at.SignedString([]byte(config.GetConfig().GetString("JWT_SECRET")))
	token, err := at.SignedString([]byte("asdasdasdasde234234234"))
	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyToken(ctx *gin.Context, bearerToken string) error {
	if bearerToken != "" && len(strings.Split(bearerToken, " ")) == 2 {
		tokenStr := strings.Split(bearerToken, " ")[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("asdasdasdasde234234234"), nil
		})
		if err != nil || !token.Valid {
			return fmt.Errorf("Invalid Token")
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Set("jwtClaims", claims)
		}
		return nil
	}
	return fmt.Errorf("Missing or Invalid Authorization header")
}

func setUserFromBearer(ctx *gin.Context, bearerToken string) error {
	if bearerToken == "" || !strings.HasPrefix(bearerToken, "Bearer ") {
		return fmt.Errorf("Missing or Invalid Authorization header")
	}
	if err := VerifyToken(ctx, bearerToken); err != nil {
		return err
	}
	jwtClaim, err := GetJWTClaim(ctx)
	if err != nil {
		return err
	}
	if jwtClaim != nil {
		claim, err := GetUserClaim(jwtClaim)
		if err != nil {
			return err
		}
		ctx.Set("user", claim)
	}
	return nil

}

func AuthRequired(required bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.Request.Header.Get("Authorization")
		if !required {
			if bearerToken == "" {
				ctx.Set("anonymous", true)
				ctx.Next()
				return
			}
			if bearerToken != "Bearer null" && strings.HasPrefix(bearerToken, "Bearer ") {
				if err := setUserFromBearer(ctx, bearerToken); err != nil {
					Error(ctx, 401, err.Error(), nil)
					ctx.Abort()
					return
				}
			}
			ctx.Next()
			return
		}

		if bearerToken != "Bearer null" && strings.HasPrefix(bearerToken, "Bearer ") {
			if err := setUserFromBearer(ctx, bearerToken); err != nil {
				Error(ctx, 401, err.Error(), nil)
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}

func GetJWTClaim(ctx *gin.Context) (jwt.MapClaims, error) {
	claims, exist := ctx.Get("jwtClaims")
	if !exist {
		Error(ctx, 401, "Unauthorized", nil)
		return nil, fmt.Errorf("Unauthorized")
	}
	return claims.(jwt.MapClaims), nil
}

func GetUserClaim(claims jwt.MapClaims) (map[string]any, error) {
	var user map[string]any
	bytes, err := json.Marshal(claims)
	if err != nil {
		return nil, nil
	}
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetKeyInCtx[T any](ctx *gin.Context, key string) (T, error) {
	anonymous, exist := ctx.Get("anonymous")
	if exist && anonymous.(bool) {
		var zero T
		return zero, nil
	}
	user, exist := ctx.Get("user")
	if !exist {
		var zero T
		return zero, fmt.Errorf("no user found in ctx")
	}
	val, ok := user.(map[string]interface{})[key]
	if !ok {
		var zero T
		return zero, fmt.Errorf("key not exist in ctx")
	}

	// Special case: convert float64 to int if needed
	var zero T
	switch any(zero).(type) {
	case int:
		if f, ok := val.(float64); ok {
			return any(int(f)).(T), nil
		}
	}
	//COnvert to real type
	typedVal, ok := val.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("value for key '%s' is not of type %T", key, zero)
	}
	return typedVal, nil
}
