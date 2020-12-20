package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/api/middleware"
	"github.com/lucthienbinh/golang_scem/internal/handler"
	"gopkg.in/validator.v2"
)

type loginRequest struct {
	Email    string `json:"email" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

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
	userAuthID, validated := handler.ValidateUserAuth(user.Email, user.Password)
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

// AppAuthRoutes for login/logout app
func AppAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/loginJSON", appLoginHandler)
	rg.GET("/logout", appLogoutHandler)
	// validate old token
	accessToken := rg.Group("/app/access-token", middleware.ValidateAppTokenForRefresh())
	accessToken.POST("/get-new", appReloginHandler)
	accessToken.GET("/check-old", appOpenHandler)
}

func appLoginHandler(c *gin.Context) {
	var user loginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	userAuthID, validated := handler.ValidateUserAuth(user.Email, user.Password)
	if validated == false {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	middleware.CreateAppToken(c, userAuthID)
	return
}

func appLogoutHandler(c *gin.Context) {
	middleware.DeleteAppToken(c)
	return
}

func appReloginHandler(c *gin.Context) {
	middleware.RefreshAppToken(c)
	return
}

func appOpenHandler(c *gin.Context) {
	middleware.CheckOldToken(c)
	return
}
