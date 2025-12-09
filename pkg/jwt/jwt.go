package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret      []byte
	expiredHour int
}

func New(secret string, expiredHour int) *JWT {
	return &JWT{
		secret:      []byte(secret),
		expiredHour: expiredHour,
	}
}

func (j *JWT) Generate(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * time.Duration(j.expiredHour)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWT) Validate(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
