package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

const SECRET_KEY = "some_secret_key_val_123123"

func GenerateJWTStaff(username string, dutyLocationId int) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	expiryTime := time.Now().Add(time.Minute * 30).Unix()

	claims["authorized"] = true
	claims["username"] = username
	claims["userType"] = "STAFF"
	claims["dutyLocationId"] = dutyLocationId
	claims["exp"] = expiryTime
	/*
	 Please note that in real world, we need to move "some_secret_key_val_123123" into something like
	 "secret.json" file of Kubernetes etc
	*/
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", 0, err
	}
	return tokenString, expiryTime, nil
}

func GenerateJWTAdmin(username string, password string) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	expiryTime := time.Now().Add(time.Minute * 30).Unix()

	claims["authorized"] = true
	claims["username"] = username
	claims["userType"] = "ADMIN"
	claims["password"] = password
	claims["exp"] = expiryTime
	/*
	 Please note that in real world, we need to move "some_secret_key_val_123123" into something like
	 "secret.json" file of Kubernetes etc
	*/
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", 0, err
	}
	return tokenString, expiryTime, nil
}

type Claims struct {
	Username       string `json:"username"`
	DutyLocationId int    `json:"dutyLocationId"`
	Password       string `json:"password"`
	UserType       string `json:"userType"`
	jwt.StandardClaims
}

func ParseToken(c *fiber.Ctx) (*Claims, error) {
	bearerHeader := c.GetReqHeaders()["Authorization"][0]
	if bearerHeader == "" {
		return nil, errors.New("missing token")
	}
	bearerToken := strings.Split(bearerHeader, " ")[1]
	tokenClaims, err := jwt.ParseWithClaims(bearerToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			timeNow := time.Now().Unix()
			fmt.Println(claims.ExpiresAt)
			fmt.Println(timeNow)
			if claims.ExpiresAt < timeNow {
				return claims, nil
			}
			return nil, errors.New("token expired")
		}
	}

	return nil, err
}
