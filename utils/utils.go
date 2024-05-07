package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

const SECRET_KEY = "some_secret_key_val_123123"

func GenerateJwtStaff(userId int, username string, dutyLocationId int) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	expiryTime := time.Now().Add(time.Minute * 30).Unix()

	claims["authorized"] = true
	claims["userId"] = userId
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

func GenerateJwtAdmin(userId int, username string, password string) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	expiryTime := time.Now().Add(time.Minute * 30).Unix()

	claims["authorized"] = true
	claims["userId"] = userId
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
	UserId         int    `json:"userId"`
	DutyLocationId int    `json:"dutyLocationId"`
	Password       string `json:"password"`
	UserType       string `json:"userType"`
	jwt.StandardClaims
}

func CtxToClaim(c *fiber.Ctx) (*Claims, string, error) {
	if c.GetReqHeaders()["Authorization"] == nil {
		return nil, "", errors.New("missing token")
	}
	if c.GetReqHeaders()["Authorization"][0] == "" {
		return nil, "", errors.New("missing token")
	}
	bearerHeader := c.GetReqHeaders()["Authorization"][0]

	if len(strings.Split(bearerHeader, " ")) != 2 {
		return nil, "", errors.New("missing token")
	}
	bearerToken := strings.Split(bearerHeader, " ")[1]

	tokenClaims, err := jwt.ParseWithClaims(bearerToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			timeNow := time.Now().Unix()
			if claims.ExpiresAt > timeNow {
				return claims, bearerToken, nil
			}
			return nil, bearerToken, errors.New("token expired")
		}
	}

	return nil, bearerToken, err
}

func ValidateJwtToken(token string) (bool, error) {

	if token == "" {
		return false, nil
	}

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			timeNow := time.Now().Unix()
			if claims.ExpiresAt > timeNow {
				return true, nil
			}
			return false, nil
		}
	}

	return false, err
}

func GetDetailsJwtToken(token string) (*Claims, error) {
	var claims *Claims
	if token == "" {
		return claims, nil
	}

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if tokenClaims != nil {
		var ok bool
		claims, ok = tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			timeNow := time.Now().Unix()
			if claims.ExpiresAt > timeNow {
				return claims, nil
			}
			return claims, nil
		}
	}

	return claims, err
}

func CtxToAuth(c *fiber.Ctx) (string, string, error) {
	authHeader := c.GetReqHeaders()["Authorization"][0]
	if authHeader == "" {
		return "", "", errors.New("missing token")
	}
	authToken := strings.Split(authHeader, " ")[1]
	decodedString, err := base64.StdEncoding.DecodeString(authToken)
	if err != nil {
		return "", "", errors.New("Decode token failed")
	}
	auth := strings.Split(string(decodedString), ":")
	return auth[0], auth[1], nil
}

func CreateMap(updateFields []string, updateValues ...interface{}) map[string]interface{} {
	var myMap = make(map[string]interface{})
	for i := 0; i < len(updateFields); i++ {
		// fmt.Printf("is nil? %s", updateValues[i] == nil)
		// if !reflect.ValueOf(updateValues[i]).IsZero() {
		myMap[updateFields[i]] = updateValues[i]
		// }
	}
	myMap["last_update_time"] = GetTimeNowString()
	return myMap
}

func JsonToMap(j interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	b, err := json.Marshal(j)
	if err != nil {
		return result, err
	}
	json.Unmarshal(b, &result)
	return result, err
}

func StringToDatetimeAddTimeZone(timeString string) (string, error) {
	var outputString string
	if timeString == "" {
		timeString = time.Time{}.Format("2006-01-02T15:04:05")
	}
	timeString = strings.Split(timeString, "+")[0]
	timeObj, err := time.Parse("2006-01-02T15:04:05", timeString)
	if err != nil {
		timeObj, err = time.Parse("200601021504", timeString)
		if err != nil {
			timeObj, err = time.Parse("20060102150405", timeString)
			if err != nil {
				timeObj, err = time.Parse("2006-01-02 15:04:05", timeString)
				if err != nil {
					timeObj, err = time.Parse("200601021504", "19700101"+timeString)
					if err != nil {
						return outputString, errors.New("Failed to translate time string (" + timeString + ") to datetime")
					}
				}
			}
		}
	}
	return timeObj.Add(8 * time.Hour).Format("2006-01-02 15:04:05"), nil
}

func GetTimeNow() time.Time {
	// all times are added +8 for the HKT time zone
	return time.Now().Add(8 * time.Hour)
}

func GetTimeNowString() string {
	return TimeInt64ToString(time.Now().Unix())
}

func TimeInt64ToString(timeInt int64) string {
	return time.Unix(timeInt, 0).Format("2006-01-02 15:04:05")
}

func StringToDatetime(timeString string) (string, error) {
	var outputString string
	if timeString == "" {
		timeString = time.Time{}.Format("2006-01-02T15:04:05")
	}
	timeString = strings.Split(timeString, "+")[0]
	timeObj, err := time.Parse("2006-01-02T15:04:05", timeString)
	if err != nil {
		timeObj, err = time.Parse("200601021504", timeString)
		if err != nil {
			timeObj, err = time.Parse("20060102150405", timeString)
			if err != nil {
				timeObj, err = time.Parse("2006-01-02 15:04:05", timeString)
				if err != nil {
					timeObj, err = time.Parse("200601021504", "19700101"+timeString)
					if err != nil {
						return outputString, errors.New("Failed to translate time string (" + timeString + ") to datetime")
					}
				}
			}
		}
	}
	return timeObj.Format("2006-01-02 15:04:05"), nil
}

func StringToRoutineResponseTime(timeString string) (string, error) {
	var outputString string
	if timeString == "" {
		timeString = time.Time{}.Format("2006-01-02T15:04:05")
	}
	timeString = strings.Split(timeString, "+")[0]
	timeObj, err := time.Parse("2006-01-02T15:04:05", timeString)
	if err != nil {
		timeObj, err = time.Parse("200601021504", timeString)
		if err != nil {
			timeObj, err = time.Parse("2006-01-02 15:04:05", timeString)
			if err != nil {
				return outputString, err
			}
		}
	}
	return timeObj.Format("1504"), nil
}

func RoutineResponseTimeToString(routineResponseTime string) (string, error) {
	var outputString string
	if routineResponseTime == "" {
		return outputString, errors.New("Empty time input")
	}
	outputString = "19700101" + routineResponseTime
	timeObj, err := time.Parse("200601021504", outputString)
	if err != nil {
		return outputString, err
	}
	return timeObj.Format("2006-01-02 15:04:05"), nil
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
