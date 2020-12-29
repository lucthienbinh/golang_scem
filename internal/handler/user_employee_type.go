package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	"gopkg.in/validator.v2"
)

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
	if err := db.First(employeeType, c.Param("id")).Error; err != nil {
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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&employeeType); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.Create(employeeType).Error; err != nil {
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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&employeeType); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
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
