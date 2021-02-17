package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/api/middleware"
	"github.com/lucthienbinh/golang_scem/internal/handler"
	"gopkg.in/validator.v2"
)

type loginRequest struct {
	Email    string `json:"email" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

// ------------------------------------- Web auth -------------------------------------

// WebAuthRoutes for login/logout web
func WebAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/loginJSON", webLoginHandler)
	rg.GET("/logout", webLogoutHandler)
}

func webLoginHandler(c *gin.Context) {
	var user loginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	_, userAuthID, validated := handler.ValidateUserAuth(user.Email, user.Password)
	if validated == false {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	middleware.CreateWebSession(c, userAuthID)
	return
}

func webLogoutHandler(c *gin.Context) {
	middleware.ClearWebSession(c)
	return
}

// ------------------------------------- App auth Redis -------------------------------------

// AppAuthRedisRoutes for login/logout app
func AppAuthRedisRoutes(rg *gin.RouterGroup) {
	rg.POST("/loginJSON", appLoginHandlerRedis)
	rg.GET("/logout", appLogoutHandlerRedis)
	// validate old token
	accessToken := rg.Group("/access-token", middleware.ValidateAppTokenForRefreshRedis())
	accessToken.POST("/get-new", appReloginHandlerRedis)
	accessToken.GET("/check-old", appOpenHandlerRedis)
}

func appLoginHandlerRedis(c *gin.Context) {
	var user loginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	frontEndContext, userAuthID, validated := handler.ValidateUserAuth(user.Email, user.Password)
	if validated == false {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	middleware.CreateAppTokenRedis(c, userAuthID, frontEndContext)
	return
}

func appLogoutHandlerRedis(c *gin.Context) {
	if err := removeFCMToken(c); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	middleware.DeleteAppTokenRedis(c)
	return
}

func appReloginHandlerRedis(c *gin.Context) {
	middleware.RefreshAppTokenRedis(c)
	return
}

func appOpenHandlerRedis(c *gin.Context) {
	middleware.CheckOldTokenRedis(c)
	return
}

// ------------------------------------- App auth BuntDB -------------------------------------

// AppAuthBuntDBRoutes for login/logout app
func AppAuthBuntDBRoutes(rg *gin.RouterGroup) {
	rg.POST("/loginJSON", appLoginHandlerBuntDB)
	rg.GET("/logout", appLogoutHandlerBuntDB)
	// validate old token
	accessToken := rg.Group("/access-token", middleware.ValidateAppTokenForRefreshBuntDB())
	accessToken.POST("/get-new", appReloginHandlerBuntDB)
	accessToken.GET("/check-old", appOpenHandlerBuntDB)
}

func appLoginHandlerBuntDB(c *gin.Context) {
	var user loginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	frontEndContext, userAuthID, validated := handler.ValidateUserAuth(user.Email, user.Password)
	if validated == false {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	middleware.CreateAppTokenBuntDB(c, userAuthID, frontEndContext)
	return
}

func appLogoutHandlerBuntDB(c *gin.Context) {
	if err := removeFCMToken(c); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	middleware.DeleteAppTokenBuntDB(c)
	return
}

func appReloginHandlerBuntDB(c *gin.Context) {
	middleware.RefreshAppTokenBuntDB(c)
	return
}

func appOpenHandlerBuntDB(c *gin.Context) {
	middleware.CheckOldTokenBuntDB(c)
	return
}

// ------------------------------------- FCM auth -------------------------------------

type appTokenRequest struct {
	Token string `json:"token" validate:"nonzero"`
}

// AppFMCToken for store FMC Token
func AppFMCToken(rg *gin.RouterGroup) {
	rg.POST("/save-token", saveFCMToken)
}

func saveFCMToken(c *gin.Context) {
	var request appTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&request); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var userAuthID uint
	var err error
	if os.Getenv("RUN_APP_AUTH") == "redis" {
		userAuthID, err = middleware.GetUserAuthIDInTokenRedis(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	} else if os.Getenv("RUN_APP_AUTH") == "buntdb" {
		userAuthID, err = middleware.GetUserAuthIDInTokenBuntDB(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
	handler.SaveFCMTokenWithUserAuthID(c, userAuthID, request.Token)
	return
}

func removeFCMToken(c *gin.Context) error {
	var userAuthID uint
	var err error
	if os.Getenv("RUN_APP_AUTH") == "redis" {
		userAuthID, err = middleware.GetUserAuthIDInTokenRedis(c)
		if err != nil {
			return err
		}
	} else if os.Getenv("RUN_APP_AUTH") == "buntdb" {
		userAuthID, err = middleware.GetUserAuthIDInTokenBuntDB(c)
		if err != nil {
			return err
		}
	}
	return handler.RemoveFCMTokenWithUserAuthID(c, userAuthID)
}
