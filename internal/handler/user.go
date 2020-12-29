package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
)

// -------------------- USER AUTHENTICATION HANDLER FUNTION --------------------

// ValidateUserAuth in database
func ValidateUserAuth(email, password string) (uint, bool) {
	userAuth := &model.UserAuthenticate{}
	if err := db.Where("email = ? AND active = true", email).First(userAuth).Error; err != nil {
		return uint(0), false
	}
	return userAuth.ID, true
}

// -------------------- FCM HANDLER FUNTION --------------------

// SaveFCMTokenWithUserAuthID to database
func SaveFCMTokenWithUserAuthID(c *gin.Context, userAuthID uint, appToken string) {
	userAuthenticate := &model.UserAuthenticate{}
	if err := db.First(userAuthenticate, userAuthID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if userAuthenticate.EmployeeID != 0 {
		employeeFCMToken := &model.EmployeeFCMToken{EmployeeID: userAuthenticate.EmployeeID, Token: appToken}
		if db.Model(&employeeFCMToken).Where("employee_id = ?", userAuthenticate.EmployeeID).Updates(&employeeFCMToken).RowsAffected == 0 {
			if err := db.Create(employeeFCMToken).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"server_response": "App token has been created!"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"server_response": "App token has been updated!"})
		return
	}
	if userAuthenticate.CustomerID != 0 {
		customerFCMToken := &model.CustomerFCMToken{CustomerID: userAuthenticate.CustomerID, Token: appToken}
		if db.Model(&customerFCMToken).Where("customer_id = ?", userAuthenticate.CustomerID).Updates(&customerFCMToken).RowsAffected == 0 {
			if err := db.Create(customerFCMToken).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"server_response": "App token has been created!"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"server_response": "App token has been updated!"})
		return
	}
}

// -------------------- COMMON FUNTION --------------------
func getIDFromParam(c *gin.Context) uint {
	rawUint64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(rawUint64)
}

func getCustomerIDByUserAuthID(userAuthID uint) (uint, error) {
	userAuthenticate := &model.UserAuthenticate{}
	if err := db.First(userAuthenticate, userAuthID).Error; err != nil {
		return uint(0), err
	}
	return userAuthenticate.CustomerID, nil
}

func getEmployeeIDByUserAuthID(userAuthID uint) (uint, error) {
	userAuthenticate := &model.UserAuthenticate{}
	if err := db.First(userAuthenticate, userAuthID).Error; err != nil {
		return uint(0), err
	}
	return userAuthenticate.CustomerID, nil
}
