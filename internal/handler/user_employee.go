package handler

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	"gopkg.in/validator.v2"
)

// -------------------- EMPLOYEE HANDLER FUNTION --------------------

// GetEmployeeListHandler in database
func GetEmployeeListHandler(c *gin.Context) {
	employeeInfoList := []model.EmployeeInfoFetchDB{}
	selectPart := "e.id, e.name, e.age, e.phone, e.gender, e.address, " +
		"e.identity_card, et.id as employee_type_id, et.name as employee_type_name, e.avatar, " +
		"dl.city as delivery_location_city, dl.district as delivery_location_district"
	leftJoin1 := "left join employee_types as et on e.employee_type_id = et.id"
	leftJoin2 := "left join delivery_locations as dl on e.delivery_location_id = dl.id"
	db.Table("employees as e").Select(selectPart).Joins(leftJoin1).Joins(leftJoin2).
		Where(" e.deleted_at is NULL ").
		Order("e.id asc").Find(&employeeInfoList)

	inputStaffTotal := 0
	deliveryStaffTotal := 0
	loadPackageStaffTotal := 0
	for i := 0; i < len(employeeInfoList); i++ {
		if employeeInfoList[i].EmployeeTypeID == uint(2) {
			inputStaffTotal++
		} else if employeeInfoList[i].EmployeeTypeID == uint(3) {
			deliveryStaffTotal++
		} else if employeeInfoList[i].EmployeeTypeID == uint(4) {
			loadPackageStaffTotal++
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"employee_list":            employeeInfoList,
		"input_staff_total":        inputStaffTotal,
		"delivery_staff_total":     deliveryStaffTotal,
		"load_package_staff_total": loadPackageStaffTotal,
	})
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
		First(employeeInfoFetchDB, c.Param("id")).Error; err != nil {
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
	newName := fmt.Sprintf("%x", b)
	createTime := fmt.Sprintf("%d", time.Now().Unix())
	newName = createTime + "_" + newName + ".jpg"
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
	if os.Getenv("SELECT_DATABASE") == "3" {
		selectPart = "dl.id as value, dl.city || ' - ' || dl.district as label "
	} else {
		selectPart = "dl.id as value, concat(dl.city, ' - ', dl.district) as label "
	}

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
	// Create employee authenticate
	if err := db.Create(userAuth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Create employee information
	employee.UserAuthID = userAuth.ID
	if err := db.Create(employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Update customer authenticate
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"server_response": "An employee has been created!"})
	return
}

// UpdateEmployeeFormData in frontend
func UpdateEmployeeFormData(c *gin.Context) {
	employee := &model.Employee{}
	if err := db.First(employee, c.Param("id")).Error; err != nil {
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
	if os.Getenv("SELECT_DATABASE") == "3" {
		selectPart = "dl.id as value, dl.city || ' - ' || dl.district as label "
	} else {
		selectPart = "dl.id as value, concat(dl.city, ' - ', dl.district) as label "
	}
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
	if err := db.First(employee, c.Param("id")).Error; err != nil {
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
	if err := db.First(employee, c.Param("id")).Error; err != nil {
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
