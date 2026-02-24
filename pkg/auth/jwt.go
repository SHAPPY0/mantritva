package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shappy0/mantritva/pkg/utils"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var issuerName = "mantritva"

type UserClaims struct {
	Name         string `json:"name"`
	UserID       uint   `json:"user_id"`
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	TokenVersion int    `json:"token_version"`
	jwt.RegisteredClaims
}

func CreateToken(claims UserClaims) (string, error) {
	if claims.UserID == 0 || claims.Name == "" || claims.UserName == "" || claims.Role == "" {
		return "", fmt.Errorf("Invalid userClaims to create token")
	}

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add((7 * 24) * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    issuerName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(ctx *gin.Context, bearerToken string) error {
	if bearerToken == "" {
		return fmt.Errorf("Missing Authorization Header")
	}
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fmt.Errorf("Invalid Authorization Token")
	}
	tokenStr := parts[1]
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected Signing Method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		},
		jwt.WithIssuer(issuerName),
		jwt.WithLeeway(5*time.Second),
	)
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("Invalid Token")
	}
	// Store claims in Gin context
	ctx.Set("jwtClaims", claims)

	return nil
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
		ctx.Set("user", jwtClaim)
	}
	return nil

}

func AuthRequired(required bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.Request.Header.Get("Authorization")
		if !required {
			if bearerToken == "" || bearerToken == "Bearer null" {
				ctx.Set("anonymous", true)
				ctx.Next()
				return
			}
			if bearerToken != "Bearer null" && strings.HasPrefix(bearerToken, "Bearer ") {
				if err := setUserFromBearer(ctx, bearerToken); err != nil {
					utils.Error(ctx, 401, err.Error(), nil)
					ctx.Abort()
					return
				}
			}
			ctx.Next()
			return
		}

		if bearerToken != "Bearer null" && strings.HasPrefix(bearerToken, "Bearer ") {
			if err := setUserFromBearer(ctx, bearerToken); err != nil {
				utils.Error(ctx, 401, err.Error(), nil)
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}

func ValidateVersion(validator TokenVersionValidator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claim, err := GetJWTClaim(ctx)
		if err != nil {
			utils.Error(ctx, 401, err.Error(), nil)
			ctx.Abort()
			return
		}
		valid, err := validator.ValidateTokenVersion(int(claim.UserID), claim.TokenVersion)
		if err != nil {
			utils.Error(ctx, 401, "Internal server error", nil)
			ctx.Abort()
			return
		}
		if !valid {
			utils.Error(ctx, 401, "Token expired, Please login again", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func GetJWTClaim(ctx *gin.Context) (*UserClaims, error) {
	claims, exist := ctx.Get("jwtClaims")
	if !exist {
		utils.Error(ctx, 401, "Unauthorized User", nil)
		return nil, fmt.Errorf("Unauthorized User")
	}
	return claims.(*UserClaims), nil
}

func GetUserClaim(claims *UserClaims) (map[string]any, error) {
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
		return zero, fmt.Errorf("No user found in ctx")
	}
	claim, err := StructToMap(user.(*UserClaims))
	if err != nil {
		var zero T
		return zero, fmt.Errorf("No user data found in ctx")
	}
	val, ok := claim[key]
	if !ok {
		var zero T
		return zero, fmt.Errorf("Key not exist in ctx")
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

func StructToMap(v interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}

	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &result)
	return result, err
}
