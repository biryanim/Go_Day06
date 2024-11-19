package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	User string
}

const TOKEN_EXP = time.Minute * 15
const SECRET_KEY = "supersecretkey"

func BuildJwtString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		User: "admin",
	})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Verify(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}
