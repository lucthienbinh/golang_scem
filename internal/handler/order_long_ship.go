package handler

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	qrcode "github.com/skip2/go-qrcode"
	"gopkg.in/validator.v2"
)

// -------------------- LONG SHIP HANDLER FUNTION --------------------

// GetLongShipListHandler in database
func GetLongShipListHandler(c *gin.Context) {
	type APILongShipList struct {
		ID                       uint   `json:"id"`
		TransportTypeID          uint   `json:"transport_type_id"`
		LicensePlate             string `json:"license_plate"`
		EstimatedTimeOfDeparture int64  `json:"estimated_time_of_departure"`
		EstimatedTimeOfArrival   int64  `json:"estimated_time_of_arrival"`
		Finished                 bool   `json:"finished"`
	}
	longShips := []APILongShipList{}
	db.Model(&model.LongShip{}).Order("id asc").Find(&longShips)

	transportTypes := []model.TransportType{}
	db.Where("same_city is ?", false).Order("id asc").Find(&transportTypes)

	c.JSON(http.StatusOK, gin.H{
		"long_ship_list":      &longShips,
		"transport_type_list": &transportTypes,
	})
	return
}

func getLongShipOrNotFound(c *gin.Context) (*model.LongShip, error) {
	longShip := &model.LongShip{}
	if err := db.First(longShip, c.Param("id")).Error; err != nil {
		return longShip, err
	}
	return longShip, nil
}

// GetLongShipHandler in database
func GetLongShipHandler(c *gin.Context) {
	longShip, err := getLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"long_ship_info": &longShip})
	return
}

// CreateLongShipFormData in frontend
func CreateLongShipFormData(c *gin.Context) {
	transportTypes := []model.TransportType{}
	db.Where("same_city is ?", false).Order("id asc").Find(&transportTypes)
	c.JSON(http.StatusOK, gin.H{"transport_type_list": &transportTypes})
	return
}

// CreateLongShipHandler in database
func CreateLongShipHandler(c *gin.Context) {
	longShip := &model.LongShip{}
	if err := c.ShouldBindJSON(&longShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&longShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Create QR code
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newName := fmt.Sprintf("%x", b)
	createTime := fmt.Sprintf("%d", time.Now().Unix())
	newName = createTime + "_" + newName + ".jpg"
	filepath := os.Getenv("QR_CODE_FILE_PATH") + newName
	if err := qrcode.WriteFile(newName, qrcode.Medium, 256, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	longShip.LSQrCode = newName
	if err := db.Create(longShip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "A long ship has been created!"})
	return
}

// UpdateLongShipFormData in frontend
func UpdateLongShipFormData(c *gin.Context) {
	longShip, err := getLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	transportTypes := []model.TransportType{}
	db.Where("same_city == ?", false).Order("id asc").Find(&transportTypes)
	c.JSON(http.StatusOK, gin.H{
		"long_ship_info":      &longShip,
		"transport_type_list": &transportTypes,
	})
	return
}

// UpdateLongShipHandler in database
func UpdateLongShipHandler(c *gin.Context) {
	longShip, err := getLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&longShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&longShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	longShip.ID = getIDFromParam(c)
	selectField := []string{"transport_type_id", "license_plate", "estimated_time_of_departure", "estimated_time_of_arrival"}
	if err = db.Model(&longShip).Select(selectField).Updates(longShip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateLSLoadPackageHandler in database
func UpdateLSLoadPackageHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	longShip, err := getLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if longShip.PackageLoaded == true {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	longShipUpdateInfo := model.LongShip{
		CurrentLocation: "Location1",
		PackageLoaded:   true,
		EmplLoadID:      employeeID,
		LoadedTime:      time.Now().Unix(),
	}
	longShip.ID = getIDFromParam(c)
	if err = db.Model(&longShip).Updates(longShipUpdateInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateLSStartVehicleHandler in database
func UpdateLSStartVehicleHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	longShip, err := getLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if longShip.PackageLoaded == false || longShip.VehicleStarted == true {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	longShipUpdateInfo := model.LongShip{
		CurrentLocation: "Location2",
		VehicleStarted:  true,
		EmplDriver1ID:   employeeID,
		StartedTime:     time.Now().Unix(),
	}
	longShip.ID = getIDFromParam(c)
	if err = db.Model(&longShip).Updates(longShipUpdateInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateLSVehicleArrivedHandler in database
func UpdateLSVehicleArrivedHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	longShip, err := getLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if longShip.VehicleStarted == false || longShip.VehicleArrived == true {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	longShipUpdateInfo := model.LongShip{
		CurrentLocation: "Location3",
		VehicleArrived:  true,
		EmplDriver2ID:   employeeID,
		ArrivedTime:     time.Now().Unix(),
	}
	longShip.ID = getIDFromParam(c)
	if err = db.Model(&longShip).Updates(longShipUpdateInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateLSUnloadPackageHandler in database
func UpdateLSUnloadPackageHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	longShip, err := getLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if longShip.VehicleArrived == false || longShip.PackageUnloaded == true {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	longShipUpdateInfo := model.LongShip{
		CurrentLocation: "Location4",
		PackageUnloaded: true,
		EmplUnloadID:    employeeID,
		UnloadedTime:    time.Now().Unix(),
		Finished:        true,
	}
	longShip.ID = getIDFromParam(c)
	if err = db.Model(&longShip).Updates(longShipUpdateInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteLongShipHandler in database
func DeleteLongShipHandler(c *gin.Context) {
	if _, err := getLongShipOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&model.LongShip{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}
