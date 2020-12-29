package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	"gopkg.in/validator.v2"
)

// -------------------- CUSTOMER HANDLER FUNTION --------------------

// GetCustomerListHandler in database
func GetCustomerListHandler(c *gin.Context) {
	customers := []model.Customer{}
	db.Order("id asc").Find(&customers)
	c.JSON(http.StatusOK, gin.H{"customer_list": &customers})
	return
}

func getCustomerOrNotFound(c *gin.Context) (*model.Customer, error) {
	customer := &model.Customer{}
	if err := db.First(customer, c.Param("id")).Error; err != nil {
		return customer, err
	}
	return customer, nil
}

// GetCustomerHandler in database
func GetCustomerHandler(c *gin.Context) {
	customer, err := getCustomerOrNotFound(c)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{"customer_info": &customer})
	return
}

// CreateCustomerHandler in database
func CreateCustomerHandler(c *gin.Context) {
	customerWithAuth := &model.CustomerWithAuth{}
	if err := c.ShouldBindJSON(&customerWithAuth); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&customerWithAuth); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	customer, userAuth := customerWithAuth.ConvertCWAToNormal()
	// Create customer authenticate
	if err := db.Create(userAuth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Create customer information
	customer.UserAuthID = userAuth.ID
	if err := db.Create(customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Update customer authenticate
	userAuth.CustomerID = customer.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Create customer credit
	customerCredit := &model.CustomerCredit{CustomerID: customer.ID, Phone: customer.Phone}
	if err := db.Create(customerCredit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "A customer has been created!"})
	return
}

// UpdateCustomerHandler in database
func UpdateCustomerHandler(c *gin.Context) {
	customer, err := getCustomerOrNotFound(c)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&customer); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	customer.ID = getIDFromParam(c)
	if err = db.Model(&customer).Omit("point").Updates(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "A customer has been updated!"})
	return
}

// UpdateCustomerCreditValidationHandler in database
func UpdateCustomerCreditValidationHandler(c *gin.Context) {
	customer, err := getCustomerOrNotFound(c)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	customerCredit := &model.CustomerCredit{}
	if err := db.Where("customer_id = ?", customer.ID).First(customerCredit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	customerCredit.ValidatePhone = true
	if err = db.Model(&customerCredit).Updates(&customerCredit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "A customer credit validation has been updated!"})
	return
}

func updateCustomerCreditBalance(accountBalance int64, customerID uint) error {
	customerCredit := &model.CustomerCredit{}
	if err := db.Where("customer_id = ?", customerID).First(customerCredit).Error; err != nil {
		return err
	}
	if err := db.Model(&customerCredit).Update("account_balance", accountBalance).Error; err != nil {
		return err
	}
	return nil
}

// UpdateCustomerCreditBalanceHandler in database
// The usecase of this function when customer report about update customer credit balance bug.
// We should we an automatic mechanism update customer credit balance when customer transfered money from bank/ online bank
func UpdateCustomerCreditBalanceHandler(c *gin.Context) {
	customer, err := getCustomerOrNotFound(c)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	accountBalanceInt64, _ := strconv.ParseInt(c.PostForm("account_balance"), 10, 64)
	if err = updateCustomerCreditBalance(accountBalanceInt64, customer.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "A customer credit balance has been updated!"})
	return
}

// DeleteCustomerHandler in database
func DeleteCustomerHandler(c *gin.Context) {
	customer, err := getCustomerOrNotFound(c)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := db.Where("customer_id = ?", customer.ID).Delete(&model.CustomerCredit{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&model.UserAuthenticate{}, customer.UserAuthID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&model.Customer{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "A customer has been deleted!"})
	return
}
