package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
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
	if err := db.First(&orderLongShip, c.Param("id")).Error; err != nil {
		return orderLongShip, err
	}
	return orderLongShip, nil
}

// CreateOrderLongShip in database
func CreateOrderLongShip(orderID uint) (uint, error) {
	orderLongShip := &model.OrderLongShip{}
	orderLongShip.OrderID = orderID
	if err := db.Create(&orderLongShip).Error; err != nil {
		return uint(0), err
	}
	return orderLongShip.ID, nil
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

// UpdateOLSLoadPackageHandler in database
func UpdateOLSLoadPackageHandler(c *gin.Context) {
	orderLongShip, err := getOrderLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	orderLongShip.ID = getIDFromParam(c)
	if err = db.Model(&orderLongShip).Updates(model.OrderLongShip{PackageLoaded: true, CurrentLocation: "Location1", EmplLoadID: 4}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOLSStartVehicleHandler in database
func UpdateOLSStartVehicleHandler(c *gin.Context) {
	orderLongShip, err := getOrderLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if orderLongShip.PackageLoaded == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	orderLongShip.ID = getIDFromParam(c)
	if err = db.Model(&orderLongShip).Updates(model.OrderLongShip{VehicleStarted: true, EmplDriverID: 3}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOLSVehicleArrivedHandler in database
func UpdateOLSVehicleArrivedHandler(c *gin.Context) {
	orderLongShip, err := getOrderLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if orderLongShip.VehicleStarted == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	orderLongShip.ID = getIDFromParam(c)
	if err = db.Model(&orderLongShip).Updates(model.OrderLongShip{VehicleArrived: true, CurrentLocation: "Location2", EmplDriverID: 3}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOLSUnloadPackageHandler in database
func UpdateOLSUnloadPackageHandler(c *gin.Context) {
	orderLongShip, err := getOrderLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if orderLongShip.VehicleArrived == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	orderLongShip.ID = getIDFromParam(c)
	if err = db.Model(&orderLongShip).Updates(model.OrderLongShip{PackageUnloaded: true, EmplDriverID: 4, Finished: true}).Error; err != nil {
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
