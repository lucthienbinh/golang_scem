package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
)

// -------------------- USER AUTHENTICATION HANDLER FUNTION --------------------

// ValidateUserAuth in database
func ValidateUserAuth(email, password string) (map[string]interface{}, uint, bool) {
	variables := make(map[string]interface{})
	userAuth := &model.UserAuthenticate{}
	if err := db.Where("email = ? AND active = true", email).First(userAuth).Error; err != nil {
		return nil, 0, false
	}
	if userAuth.EmployeeID != 0 {
		variables["employee_id"] = userAuth.EmployeeID
		variables["customer_id"] = 0
		employee := &model.Employee{}
		if err := db.First(employee, userAuth.EmployeeID).Error; err != nil {
			return nil, 0, false
		}
		variables["name"] = employee.Name
		variables["address"] = employee.Address
		variables["phone"] = employee.Phone
		variables["gender"] = employee.Gender
		variables["age"] = employee.Age
	}
	if userAuth.CustomerID != 0 {
		variables["customer_id"] = userAuth.CustomerID
		variables["employee_id"] = 0
		customer := &model.Customer{}
		if err := db.First(customer, userAuth.CustomerID).Error; err != nil {
			return nil, 0, false
		}
		variables["name"] = customer.Name
		variables["address"] = customer.Address
		variables["phone"] = customer.Phone
		variables["gender"] = customer.Gender
		variables["age"] = customer.Age
	}
	return variables, userAuth.ID, true
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
		employeeFCMToken := &model.UserFCMToken{EmployeeID: userAuthenticate.EmployeeID, Token: appToken}
		if db.Model(employeeFCMToken).Where("employee_id = ?", userAuthenticate.EmployeeID).Updates(employeeFCMToken).RowsAffected == 0 {
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
		customerFCMToken := &model.UserFCMToken{CustomerID: userAuthenticate.CustomerID, Token: appToken}
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

// RemoveFCMTokenWithUserAuthID to database
func RemoveFCMTokenWithUserAuthID(c *gin.Context, userAuthID uint) error {
	userAuthenticate := &model.UserAuthenticate{}
	if err := db.First(userAuthenticate, userAuthID).Error; err != nil {
		return err
	}
	if userAuthenticate.EmployeeID != 0 {
		employeeFCMToken := &model.UserFCMToken{}
		if err := db.Model(employeeFCMToken).Where("employee_id = ?", userAuthenticate.EmployeeID).Delete(employeeFCMToken).Error; err != nil {
			return err
		}

	}
	if userAuthenticate.CustomerID != 0 {
		customerFCMToken := &model.UserFCMToken{}
		if err := db.Model(customerFCMToken).Where("customer_id = ?", userAuthenticate.CustomerID).Delete(customerFCMToken).Error; err != nil {
			return err
		}
	}
	return nil
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
	return userAuthenticate.EmployeeID, nil
}
