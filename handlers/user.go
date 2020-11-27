package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/middlewares"
	"github.com/lucthienbinh/golang_scem/models"
)

// LoginHandler check user information
func LoginHandler(c *gin.Context) {
	var json models.Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if json.Email == "" || json.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing field"})
		return
	}
	userAuth := &models.UserAuthenticate{}
	if err := db.Where("email = ? AND active = true", json.Email).First(&userAuth).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	middlewares.CreateSession(c, userAuth)
	return
}

// LogoutHandler remove user session
func LogoutHandler(c *gin.Context) {
	middlewares.ClearSession(c)
	return
}

// -------------------- CUSTOMER HANDLER FUNTION --------------------

// GetCustomerListHandler in database
func GetCustomerListHandler(c *gin.Context) {
	customers := []models.Customer{}
	db.Find(&customers)
	c.JSON(http.StatusOK, gin.H{"customer_list": &customers})
	return
}

func getCustomerOrNotFound(c *gin.Context) (*models.Customer, error) {
	customer := &models.Customer{}
	if err := db.First(&customer, c.Param("id")).Error; err != nil {
		return customer, err
	}
	return customer, nil
}

// GetCustomerHandler in database
func GetCustomerHandler(c *gin.Context) {
	customer, err := getCustomerOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"customer_info": &customer})
	return
}

// CreateCustomerHandler in database
func CreateCustomerHandler(c *gin.Context) {
	customerWithAuth := &models.CustomerWithAuth{}
	if err := c.ShouldBindJSON(&customerWithAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer, userAuth := customerWithAuth.ConvertCWAToNormal()
	if err := db.Create(&userAuth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if userAuth.ID != 0 {
		customer.UserAuthID = userAuth.ID
		if err := db.Create(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "A customer has been created!"})
	return
}

// UpdateCustomerHandler in database
func UpdateCustomerHandler(c *gin.Context) {
	customer, err := getCustomerOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fetchCustomer := &models.Customer{}
	if err := c.ShouldBindJSON(&fetchCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// dont update these field
	fetchCustomer.ID = customer.ID
	fetchCustomer.Point = customer.Point
	if err = db.Model(&fetchCustomer).Updates(&fetchCustomer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteCustomerHandler in database
func DeleteCustomerHandler(c *gin.Context) {
	if _, err := getCustomerOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&models.Customer{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}

// -------------------- EMPLOYEE HANDLER FUNTION --------------------

// GetEmployeeListHandler in database
func GetEmployeeListHandler(c *gin.Context) {
	employees := []models.Employee{}
	db.Find(&employees)
	c.JSON(http.StatusOK, gin.H{"employee_list": &employees})
	return
}

func getEmployeeOrNotFound(c *gin.Context) (*models.Employee, error) {
	employee := &models.Employee{}
	if err := db.First(&employee, c.Param("id")).Error; err != nil {
		return employee, err
	}
	return employee, nil
}

// GetEmployeeHandler in database
func GetEmployeeHandler(c *gin.Context) {
	employee, err := getEmployeeOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"employee_info": &employee})
	return
}

// CreateEmployeeHandler in database
func CreateEmployeeHandler(c *gin.Context) {
	employeeWithAuth := &models.EmployeeWithAuth{}
	if err := c.ShouldBindJSON(&employeeWithAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	employee, userAuth := employeeWithAuth.ConvertEWAToNormal()
	if err := db.Create(&userAuth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if userAuth.ID != 0 {
		employee.UserAuthID = userAuth.ID
		if err := db.Create(&employee).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"server_response": "An employee has been created!"})
	return
}

// UpdateEmployeeHandler in database
func UpdateEmployeeHandler(c *gin.Context) {
	employee, err := getEmployeeOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fetchEmployee := &models.Employee{}
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// dont update these field
	fetchEmployee.ID = employee.ID
	if err = db.Model(&fetchEmployee).Updates(&fetchEmployee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteEmployeeHandler in database
func DeleteEmployeeHandler(c *gin.Context) {
	if _, err := getEmployeeOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&models.Employee{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}

// -------------------- DELIVERY LOCATION HANDLER FUNTION --------------------

// GetDeliveryLocationListHandler in database
func GetDeliveryLocationListHandler(c *gin.Context) {
	deliveryLocations := []models.DeliveryLocation{}
	db.Find(&deliveryLocations)
	c.JSON(http.StatusOK, gin.H{"delivery_type_list": &deliveryLocations})
	return
}

func getDeliveryLocationOrNotFound(c *gin.Context) (*models.DeliveryLocation, error) {
	deliveryLocation := &models.DeliveryLocation{}
	if err := db.First(&deliveryLocation, c.Param("id")).Error; err != nil {
		return deliveryLocation, err
	}
	return deliveryLocation, nil
}

// GetDeliveryLocationHandler in database
func GetDeliveryLocationHandler(c *gin.Context) {
	deliveryLocation, err := getDeliveryLocationOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"delivery_type_info": &deliveryLocation})
	return
}

// CreateDeliveryLocationHandler in database
func CreateDeliveryLocationHandler(c *gin.Context) {
	deliveryLocation := &models.DeliveryLocation{}
	if err := c.ShouldBindJSON(&deliveryLocation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&deliveryLocation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "An employee type has been created!"})
	return
}

// UpdateDeliveryLocationHandler in database
func UpdateDeliveryLocationHandler(c *gin.Context) {
	deliveryLocation, err := getDeliveryLocationOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&deliveryLocation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = db.Model(&deliveryLocation).Updates(&deliveryLocation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteDeliveryLocationHandler in database
func DeleteDeliveryLocationHandler(c *gin.Context) {
	if _, err := getDeliveryLocationOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&models.DeliveryLocation{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}

// -------------------- EMPLOYEE TYPE HANDLER FUNTION --------------------

// GetEmployeeTypeListHandler in database
func GetEmployeeTypeListHandler(c *gin.Context) {
	employeeTypes := []models.EmployeeType{}
	db.Find(&employeeTypes)
	c.JSON(http.StatusOK, gin.H{"employee_type_list": &employeeTypes})
	return
}

func getEmployeeTypeOrNotFound(c *gin.Context) (*models.EmployeeType, error) {
	employeeType := &models.EmployeeType{}
	if err := db.First(&employeeType, c.Param("id")).Error; err != nil {
		return employeeType, err
	}
	return employeeType, nil
}

// GetEmployeeTypeHandler in database
func GetEmployeeTypeHandler(c *gin.Context) {
	employeeType, err := getEmployeeTypeOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"employee_type_info": &employeeType})
	return
}

// CreateEmployeeTypeHandler in database
func CreateEmployeeTypeHandler(c *gin.Context) {
	employeeType := &models.EmployeeType{}
	if err := c.ShouldBindJSON(&employeeType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&employeeType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "An employee type has been created!"})
	return
}

// UpdateEmployeeTypeHandler in database
func UpdateEmployeeTypeHandler(c *gin.Context) {
	employeeType, err := getEmployeeTypeOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&employeeType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = db.Model(&employeeType).Updates(&employeeType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteEmployeeTypeHandler in database
func DeleteEmployeeTypeHandler(c *gin.Context) {
	if _, err := getEmployeeTypeOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&models.EmployeeType{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}
