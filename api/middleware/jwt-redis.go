package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
)

// Source: https://learn.vonage.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr

var client *redis.Client

// RunAppAuthRedis to connect redis server
func RunAppAuthRedis() error {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

// -------------------- Public function --------------------

// CreateAppTokenRedis after logged in successful
func CreateAppTokenRedis(c *gin.Context, userAuthID uint, frontEndContext map[string]interface{}) {

	ts, err := createToken(userAuthID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := createAuth(userAuthID, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
		"token_type":    "Bearer ",
		"expires":       strconv.FormatInt(ts.AtExpires, 10),
	}
	c.JSON(http.StatusCreated, gin.H{
		"tokens":            tokens,
		"front_end_context": frontEndContext,
	})
	return
}

// ValidateAppTokenRedis secure private routes
func ValidateAppTokenRedis() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tokenValid(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": err.Error()})
		} else {
			c.Next()
		}

	}
}

// ValidateAppTokenForRefreshRedis secure private routes
func ValidateAppTokenForRefreshRedis() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tokenValid(c.Request)
		if err != nil {
			if (err.Error() == "Token is expired") && (c.FullPath() == "/app-auth/access-token/get-new") {
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": err.Error()})
			}

		}
		if err == nil {
			c.Next()
		}
	}
}

// DeleteAppTokenRedis after logged out successful
func DeleteAppTokenRedis(c *gin.Context) {
	metadata, err := extractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	delErr := deleteTokens(metadata)
	if delErr != nil {
		c.JSON(http.StatusUnauthorized, delErr.Error())
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

// CheckOldTokenRedis to check when open app
func CheckOldTokenRedis(c *gin.Context) {
	c.JSON(http.StatusOK, "Valid app access token")
}

// RefreshAppTokenRedis function
func RefreshAppTokenRedis(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		fmt.Println("the error: ", err)
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		userIDConvert := uint(userID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := deleteAuth(refreshUUID)
		if delErr != nil || deleted == 0 { //if any goes wrong
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := createToken(userIDConvert)
		if createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		saveErr := createAuth(userIDConvert, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
			"token_type":    "Bearer ",
			"expires":       strconv.FormatInt(ts.AtExpires, 10),
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}

// GetUserAuthIDInTokenRedis function
func GetUserAuthIDInTokenRedis(c *gin.Context) (uint, error) {
	//Extract the access token metadata
	metadata, err := extractTokenMetadata(c.Request)
	if err != nil {
		return uint(0), err
	}
	userID, err := fetchAuth(metadata)
	return uint(userID), nil
}

// -------------------- Private function --------------------

type accessDetails struct {
	AccessUUID string
	UserID     uint
}

type tokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

func createToken(userid uint) (*tokenDetails, error) {
	td := &tokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = td.AccessUUID + "++" + strconv.Itoa(int(userid))

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func createAuth(userid uint, td *tokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUUID, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUUID, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	if tokenString == "" {
		return nil, errors.New("unauthorized")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func tokenValid(r *http.Request) error {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return err
	}
	return nil
}

func extractTokenMetadata(r *http.Request) (*accessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		userIDConvert := uint(userID)
		return &accessDetails{
			AccessUUID: accessUUID,
			UserID:     userIDConvert,
		}, nil
	}
	return nil, err
}

func deleteAuth(givenUUID string) (int64, error) {
	deleted, err := client.Del(givenUUID).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func deleteTokens(authD *accessDetails) error {
	//get the refresh uuid
	refreshUUID := fmt.Sprintf("%s++%d", authD.AccessUUID, authD.UserID)
	//delete access token
	deletedAt, err := client.Del(authD.AccessUUID).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := client.Del(refreshUUID).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

// -------------------- Author example --------------------
// Fetch data in redis with extractTokenMetadata()

func fetchAuth(authD *accessDetails) (uint, error) {
	userid, err := client.Get(authD.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	userIDConvert := uint(userID)
	if authD.UserID != userIDConvert {
		return 0, errors.New("unauthorized")
	}
	return userIDConvert, nil
}

func createTodo(c *gin.Context) {

	//Extract the access token metadata
	metadata, err := extractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := fetchAuth(metadata)
	userIDConvert := uint(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	// you can proceed to save the Todo to a database
	// but we will just return it to the caller:

	c.JSON(http.StatusCreated, userIDConvert)
}
