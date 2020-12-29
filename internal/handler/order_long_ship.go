package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	"gopkg.in/validator.v2"
)

// -------------------- ORDER LONG SHIP HANDLER FUNTION --------------------

// GetOrderLongShipListHandler in database
func GetOrderLongShipListHandler(c *gin.Context) {
	orderLongShips := []model.OrderLongShip{}
	db.Order("id asc").Find(&orderLongShips)
	c.JSON(http.StatusOK, gin.H{"order_long_ship_list": &orderLongShips})
	return
}

func getOrderLongShipOrNotFound(c *gin.Context) (*model.OrderLongShip, error) {
	orderLongShip := &model.OrderLongShip{}
	if err := db.First(orderLongShip, c.Param("id")).Error; err != nil {
		return orderLongShip, err
	}
	return orderLongShip, nil
}

// GetOrderLongShipHandler in database
func GetOrderLongShipHandler(c *gin.Context) {
	orderLongShip, err := getOrderLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order_long_ship_info": &orderLongShip})
	return
}

// CreateLongShipFormData in frontend
func CreateLongShipFormData(c *gin.Context) {
	transportTypes := []model.TransportType{}
	db.Where("same_city == ?", false).Order("id asc").Find(&transportTypes)
	c.JSON(http.StatusOK, gin.H{"transport_type_list": &transportTypes})
	return
}

// CreateOrderLongShipHandler in database
func CreateOrderLongShipHandler(c *gin.Context) {
	orderLongShip := &model.OrderLongShip{}
	if err := c.ShouldBindJSON(&orderLongShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&orderLongShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.Create(orderLongShip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "An order long ship has been created!"})
	return
}

// UpdateLongShipFormData in frontend
func UpdateLongShipFormData(c *gin.Context) {
	orderLongShip, err := getOrderLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	transportTypes := []model.TransportType{}
	db.Where("same_city == ?", false).Order("id asc").Find(&transportTypes)
	c.JSON(http.StatusOK, gin.H{
		"order_long_ship_info": &orderLongShip,
		"transport_type_list":  &transportTypes,
	})
	return
}

// UpdateOrderLongShipHandler in database
func UpdateOrderLongShipHandler(c *gin.Context) {
	orderLongShip, err := getOrderLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&orderLongShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&orderLongShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	orderLongShip.ID = getIDFromParam(c)
	selectField := []string{"transport_type_id", "license_plate", "estimated_time_of_departure", "estimated_time_of_arrival"}
	if err = db.Model(&orderLongShip).Select(selectField).Updates(orderLongShip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOLSLoadPackageHandler in database
func UpdateOLSLoadPackageHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	orderLongShip, err := getOrderLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	orderLongShipUpdateInfo := model.OrderLongShip{
		CurrentLocation: "Location1",
		PackageLoaded:   true,
		EmplLoadID:      employeeID,
		LoadedTime:      time.Now().Unix(),
	}
	orderLongShip.ID = getIDFromParam(c)
	if err = db.Model(&orderLongShip).Updates(orderLongShipUpdateInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOLSStartVehicleHandler in database
func UpdateOLSStartVehicleHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	orderLongShip, err := getOrderLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if orderLongShip.PackageLoaded == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	orderLongShipUpdateInfo := model.OrderLongShip{
		CurrentLocation: "Location2",
		VehicleStarted:  true,
		EmplDriver1ID:   employeeID,
		StartedTime:     time.Now().Unix(),
	}
	orderLongShip.ID = getIDFromParam(c)
	if err = db.Model(&orderLongShip).Updates(orderLongShipUpdateInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOLSVehicleArrivedHandler in database
func UpdateOLSVehicleArrivedHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	orderLongShip, err := getOrderLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if orderLongShip.VehicleStarted == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	orderLongShipUpdateInfo := model.OrderLongShip{
		CurrentLocation: "Location3",
		VehicleArrived:  true,
		EmplDriver2ID:   employeeID,
		ArrivedTime:     time.Now().Unix(),
	}
	orderLongShip.ID = getIDFromParam(c)
	if err = db.Model(&orderLongShip).Updates(orderLongShipUpdateInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOLSUnloadPackageHandler in database
func UpdateOLSUnloadPackageHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	orderLongShip, err := getOrderLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if orderLongShip.VehicleArrived == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	orderLongShipUpdateInfo := model.OrderLongShip{
		CurrentLocation: "Location3",
		PackageUnloaded: true,
		EmplUnloadID:    employeeID,
		UnloadedTime:    time.Now().Unix(),
	}
	orderLongShip.ID = getIDFromParam(c)
	if err = db.Model(&orderLongShip).Updates(orderLongShipUpdateInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteOrderLongShipHandler in database
func DeleteOrderLongShipHandler(c *gin.Context) {
	if _, err := getOrderLongShipOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&model.OrderLongShip{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}
