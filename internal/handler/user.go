package handler

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	"gopkg.in/validator.v2"
)

// -------------------- USER AUTHENTICATION HANDLER FUNTION --------------------

// ValidateUserAuth in database
func ValidateUserAuth(email, password string) (uint, bool) {
	userAuth := &model.UserAuthenticate{}
	if err := db.Where("email = ? AND active = true", email).First(&userAuth).Error; err != nil {
		return uint(0), false
	}
	return userAuth.ID, true
}

// -------------------- FCM HANDLER FUNTION --------------------

// SaveFCMTokenWithUserAuthID to database
func SaveFCMTokenWithUserAuthID(c *gin.Context, userAuthID uint, appToken string) {
	userAuthenticate := &model.UserAuthenticate{}
	if err := db.First(&userAuthenticate, userAuthID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if userAuthenticate.UserType == 1 {
		employee := &model.Employee{}
		if err := db.Where("user_auth_id = ?", userAuthID).First(&employee).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		employeeFCMToken := &model.EmployeeFCMToken{}
		employeeFCMToken.EmployeeID = employee.ID
		employeeFCMToken.Token = appToken
		if db.Model(&employeeFCMToken).Where("employee_id = ?", employee.ID).Updates(&employeeFCMToken).RowsAffected == 0 {
			if err := db.Create(&employeeFCMToken).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"server_response": "App token has been created!"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"server_response": "App token has been updated!"})
		return
	}
	if userAuthenticate.UserType == 2 {
		customer := &model.Customer{}
		if err := db.Where("user_auth_id = ?", userAuthID).First(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		customerFCMToken := &model.CustomerFCMToken{}
		customerFCMToken.CustomerID = customer.ID
		customerFCMToken.Token = appToken
		if db.Model(&customerFCMToken).Where("customer_id = ?", customer.ID).Updates(&customerFCMToken).RowsAffected == 0 {
			if err := db.Create(&customerFCMToken).Error; err != nil {
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
	rawUint64, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	return uint(rawUint64)
}

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
	if err := db.First(&customer, c.Param("id")).Error; err != nil {
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
	userAuth.UserType = 2
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
	customer2 := &model.Customer{}
	customer2.ID = 2
	if err = db.Model(&customer2).Omit("point").Updates(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "A customer has been updated!"})
	return
}

// DeleteCustomerHandler in database
func DeleteCustomerHandler(c *gin.Context) {
	customer, err := getCustomerOrNotFound(c)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
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

// -------------------- EMPLOYEE HANDLER FUNTION --------------------

// GetEmployeeListHandler in database
func GetEmployeeListHandler(c *gin.Context) {
	employeeInfoList := []model.EmployeeInfoFetchDB{}
	selectPart := "e.id, e.name, e.age, e.phone, e.gender, e.address, " +
		"e.identity_card, et.name as employee_type_name, e.avatar, " +
		"dl.city as delivery_location_city, dl.district as delivery_location_district"
	leftJoin1 := "left join employee_types as et on e.employee_type_id = et.id"
	leftJoin2 := "left join delivery_locations as dl on e.delivery_location_id = dl.id"
	db.Table("employees as e").Select(selectPart).Joins(leftJoin1).Joins(leftJoin2).
		Where(" e.deleted_at is NULL ").
		Order("e.id asc").Find(&employeeInfoList)
	c.JSON(http.StatusOK, gin.H{"employee_list": employeeInfoList})
	return
}

// GetEmployeeHandler in database
func GetEmployeeHandler(c *gin.Context) {

	employeeInfoFetchDB := &model.EmployeeInfoFetchDB{}
	selectPart := "e.id, e.name, e.age, e.phone, e.gender, e.address, " +
		"e.identity_card, et.name as employee_type_name, e.avatar, " +
		"dl.city as delivery_location_city, dl.district as delivery_location_district"
	leftJoin1 := "left join employee_types as et on e.employee_type_id = et.id"
	leftJoin2 := "left join delivery_locations as dl on e.delivery_location_id = dl.id"
	if err := db.Table("employees as e").Select(selectPart).Joins(leftJoin1).Joins(leftJoin2).
		Where(" e.deleted_at is NULL ").
		First(&employeeInfoFetchDB, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{"employee_info": &employeeInfoFetchDB})
	return
}

// ImageEmployeeHandler updload image of employee
func ImageEmployeeHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Println(file.Filename)

	b := make([]byte, 8)

	if _, err := rand.Read(b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	extension := strings.Split(file.Filename, ".")
	newName := fmt.Sprintf("%x", b)
	createTime := fmt.Sprintf("%d", time.Now().Unix())
	newName = createTime + "_" + newName + "." + extension[1]
	filepath := os.Getenv("IMAGE_FILE_PATH") + newName

	// Upload the file to specific dst.
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"filename": newName})
	return
}

// CreateEmployeeFormData in frontend
func CreateEmployeeFormData(c *gin.Context) {
	employeeTypeOptions := []model.SelectStuct{}
	selectPart := "et.id as value, et.name as label "
	db.Table("employee_types as et").Select(selectPart).Order("et.id asc").Find(&employeeTypeOptions)
	for i := 0; i < len(employeeTypeOptions); i++ {
		employeeTypeOptions[i].Name = "employee_type_id"
	}

	deliveryLocationOptions := []model.SelectStuct{}
	selectPart = "dl.id as value, concat(dl.city, ' - ', dl.district) as label "
	db.Table("delivery_locations as dl").Select(selectPart).Order("dl.id asc").Find(&deliveryLocationOptions)
	for i := 0; i < len(deliveryLocationOptions); i++ {
		deliveryLocationOptions[i].Name = "delivery_location_id"
	}

	c.JSON(http.StatusOK, gin.H{
		"et_options": &employeeTypeOptions,
		"dl_options": &deliveryLocationOptions,
	})
	return
}

// CreateEmployeeHandler in database
func CreateEmployeeHandler(c *gin.Context) {
	employeeWithAuth := &model.EmployeeWithAuth{}
	if err := c.ShouldBindJSON(&employeeWithAuth); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&employeeWithAuth); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	employee, userAuth := employeeWithAuth.ConvertEWAToNormal()
	userAuth.UserType = 1
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

	c.JSON(http.StatusCreated, gin.H{"server_response": "An employee has been created!"})
	return
}

// UpdateEmployeeFormData in frontend
func UpdateEmployeeFormData(c *gin.Context) {
	employee := &model.Employee{}
	if err := db.First(&employee, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	employeeTypeOptions := []model.SelectStuct{}
	selectPart := "et.id as value, et.name as label "
	db.Table("employee_types as et").Select(selectPart).Order("et.id asc").Find(&employeeTypeOptions)
	for i := 0; i < len(employeeTypeOptions); i++ {
		employeeTypeOptions[i].Name = "employee_type_id"
	}

	deliveryLocationOptions := []model.SelectStuct{}
	selectPart = "dl.id as value, concat(dl.city, ' - ', dl.district) as label "
	db.Table("delivery_locations as dl").Select(selectPart).Order("dl.id asc").Find(&deliveryLocationOptions)
	for i := 0; i < len(deliveryLocationOptions); i++ {
		deliveryLocationOptions[i].Name = "delivery_location_id"
	}

	c.JSON(http.StatusOK, gin.H{
		"et_options":    &employeeTypeOptions,
		"dl_options":    &deliveryLocationOptions,
		"employee_info": &employee,
	})
	return
}

// UpdateEmployeeHandler in database
func UpdateEmployeeHandler(c *gin.Context) {
	employee := &model.Employee{}
	if err := db.First(&employee, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&employee); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	employee.ID = getIDFromParam(c)
	if err := db.Model(&employee).Updates(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteEmployeeHandler in database
func DeleteEmployeeHandler(c *gin.Context) {
	employee := &model.Employee{}
	if err := db.First(&employee, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := db.Delete(&model.UserAuthenticate{}, employee.UserAuthID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Delete(&model.Employee{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}

// -------------------- DELIVERY LOCATION HANDLER FUNTION --------------------

// GetDeliveryLocationListHandler in database
func GetDeliveryLocationListHandler(c *gin.Context) {
	deliveryLocations := []model.DeliveryLocation{}
	db.Order("id asc").Find(&deliveryLocations)
	c.JSON(http.StatusOK, gin.H{"delivery_location_list": &deliveryLocations})
	return
}

func getDeliveryLocationOrNotFound(c *gin.Context) (*model.DeliveryLocation, error) {
	deliveryLocation := &model.DeliveryLocation{}
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
	c.JSON(http.StatusOK, gin.H{"delivery_location_info": &deliveryLocation})
	return
}

// CreateDeliveryLocationHandler in database
func CreateDeliveryLocationHandler(c *gin.Context) {
	deliveryLocation := &model.DeliveryLocation{}
	if err := c.ShouldBindJSON(&deliveryLocation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate(&deliveryLocation); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.Create(&deliveryLocation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "A delivery location has been created!"})
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
	if err := validator.Validate(&deliveryLocation); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	deliveryLocation.ID = getIDFromParam(c)
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
	if err := db.Delete(&model.DeliveryLocation{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}

// -------------------- EMPLOYEE TYPE HANDLER FUNTION --------------------

// GetEmployeeTypeListHandler in database
func GetEmployeeTypeListHandler(c *gin.Context) {
	employeeTypes := []model.EmployeeType{}
	db.Order("id asc").Find(&employeeTypes)
	c.JSON(http.StatusOK, gin.H{"employee_type_list": &employeeTypes})
	return
}

func getEmployeeTypeOrNotFound(c *gin.Context) (*model.EmployeeType, error) {
	employeeType := &model.EmployeeType{}
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
	employeeType := &model.EmployeeType{}
	if err := c.ShouldBindJSON(&employeeType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&employeeType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "An employee type has been created!"})
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
	employeeType.ID = getIDFromParam(c)
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
	if err := db.Delete(&model.EmployeeType{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}
