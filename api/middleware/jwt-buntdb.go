package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
)

// Source: https://learn.vonage.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr
// Source: https://github.com/tidwall/buntdb

// -------------------- Public function --------------------

// CreateAppTokenBuntDB after logged in successful
func CreateAppTokenBuntDB(c *gin.Context, userAuthID uint, frontEndContext map[string]interface{}) {

	ts, err := createToken(userAuthID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := createAuthBuntDB(userAuthID, ts)
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

// ValidateAppTokenBuntDB secure private routes
func ValidateAppTokenBuntDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tokenValid(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": err.Error()})
		} else {
			c.Next()
		}

	}
}

// ValidateAppTokenForRefreshBuntDB secure private routes
func ValidateAppTokenForRefreshBuntDB() gin.HandlerFunc {
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

// DeleteAppTokenBuntDB after logged out successful
func DeleteAppTokenBuntDB(c *gin.Context) {
	metadata, err := extractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	delErr := deleteTokensBuntDB(metadata)
	if delErr != nil {
		c.JSON(http.StatusUnauthorized, delErr.Error())
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

// CheckOldTokenBuntDB to check when open app
func CheckOldTokenBuntDB(c *gin.Context) {
	c.JSON(http.StatusOK, "Valid app access token")
}

// RefreshAppTokenBuntDB function
func RefreshAppTokenBuntDB(c *gin.Context) {
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
		delErr := deleteAuthBuntDB(refreshUUID)
		if delErr != nil { //if any goes wrong
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
		saveErr := createAuthBuntDB(userIDConvert, ts)
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

// GetUserAuthIDInTokenBuntDB function
func GetUserAuthIDInTokenBuntDB(c *gin.Context) (uint, error) {
	//Extract the access token metadata
	metadata, err := extractTokenMetadata(c.Request)
	if err != nil {
		return uint(0), err
	}
	userID, err := fetchAuthBuntDB(metadata)
	return uint(userID), nil
}

// -------------------- Private function --------------------

func fetchAuthBuntDB(authD *accessDetails) (uint, error) {
	var userid string
	db, err := buntdb.Open(os.Getenv("BUNTDB_DSN"))
	if err != nil {
		return 0, errors.New("unauthorized")
	}
	defer db.Close()

	err = db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(authD.AccessUUID)
		if err != nil {
			return err
		}
		userid = val
		return nil
	})
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

func createAuthBuntDB(userid uint, td *tokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	db, err := buntdb.Open(os.Getenv("BUNTDB_DSN"))
	if err != nil {
		return err
	}
	defer db.Close()

	errUpdate := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(td.AccessUUID, strconv.Itoa(int(userid)), &buntdb.SetOptions{Expires: true, TTL: at.Sub(now)})
		return err
	})
	if errUpdate != nil {
		return errUpdate
	}

	errRefresh := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(td.RefreshUUID, strconv.Itoa(int(userid)), &buntdb.SetOptions{Expires: true, TTL: rt.Sub(now)})
		return err
	})
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func deleteAuthBuntDB(givenUUID string) error {

	db, err := buntdb.Open(os.Getenv("BUNTDB_DSN"))
	if err != nil {
		return err
	}
	defer db.Close()

	errDelete := db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(givenUUID)
		return err
	})
	if errDelete != nil {
		return errDelete
	}

	return nil
}

func deleteTokensBuntDB(authD *accessDetails) error {
	//get the refresh uuid
	refreshUUID := fmt.Sprintf("%s++%d", authD.AccessUUID, authD.UserID)
	// Connect DB
	db, err := buntdb.Open(os.Getenv("BUNTDB_DSN"))
	if err != nil {
		return err
	}
	defer db.Close()

	//delete access token
	errDeleteAt := db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(authD.AccessUUID)
		return err
	})
	if errDeleteAt != nil {
		return errDeleteAt
	}

	//delete refresh token
	errDeleteRt := db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(refreshUUID)
		return err
	})
	if errDeleteRt != nil {
		return errDeleteRt
	}

	return nil
}
