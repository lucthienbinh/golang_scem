package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
)

// -------------------- ORDER SHORT SHIP HANDLER FUNTION --------------------

// GetOrderShortShipListHandler in database
func GetOrderShortShipListHandler(c *gin.Context) {
	orderShortShips := []model.OrderShortShip{}
	db.Order("id asc").Find(&orderShortShips)
	c.JSON(http.StatusOK, gin.H{"order_long_ship_list": &orderShortShips})
	return
}

func getOrderShortShipOrNotFound(c *gin.Context) (*model.OrderShortShip, error) {
	orderShortShip := &model.OrderShortShip{}
	if err := db.First(orderShortShip, c.Param("id")).Error; err != nil {
		return orderShortShip, err
	}
	return orderShortShip, nil
}

// GetOrderShortShipHandler in database
func GetOrderShortShipHandler(c *gin.Context) {
	orderShortShip, err := getOrderShortShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order_long_ship_info": &orderShortShip})
	return
}

// UpdateOSSShipperReceivedHandler in database
func UpdateOSSShipperReceivedHandler(c *gin.Context) {
	orderShortShip, err := getOrderShortShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	orderShortShip.ID = getIDFromParam(c)
	if err = db.Model(&orderShortShip).Updates(model.OrderShortShip{ShipperReceived: true}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOSSShipperCalledHandler in database
func UpdateOSSShipperCalledHandler(c *gin.Context) {
	orderShortShip, err := getOrderShortShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if orderShortShip.ShipperReceived == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	timeConfirmed := c.PostForm("time_confirmed")
	timeConfirmedParsed, _ := strconv.ParseInt(timeConfirmed, 10, 64)

	orderShortShip.ID = getIDFromParam(c)
	if err = db.Model(&orderShortShip).Updates(model.OrderShortShip{ShipperCalled: true, TimeConfirmed: timeConfirmedParsed}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOSSShipperShippedHandler in database
func UpdateOSSShipperShippedHandler(c *gin.Context) {
	orderShortShip, err := getOrderShortShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if orderShortShip.ShipperCalled == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	orderShortShip.ID = getIDFromParam(c)
	if err = db.Model(&orderShortShip).Updates(model.OrderShortShip{ShipperShipped: true}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOSSCusReceiveConfirmedHandler in database
func UpdateOSSCusReceiveConfirmedHandler(c *gin.Context) {
	orderShortShip, err := getOrderShortShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if orderShortShip.ShipperShipped == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	orderShortShip.ID = getIDFromParam(c)
	if err = db.Model(&orderShortShip).Updates(model.OrderShortShip{CusReceiveConfirmed: true, Finished: true}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOSSShipperConfirmedHandler in database
func UpdateOSSShipperConfirmedHandler(c *gin.Context) {
	orderShortShip, err := getOrderShortShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if orderShortShip.ShipperShipped == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	b := make([]byte, 8)

	if _, err := rand.Read(b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	extension := strings.Split(file.Filename, ".")
	newName := fmt.Sprintf("%x", b)
	createTime := fmt.Sprintf("%d", time.Now().Unix())
	newName = createTime + "_" + newName + "." + extension[1]
	filepath := "public/upload/images/" + newName

	// Upload the file to specific dst.
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderShortShip.ID = getIDFromParam(c)
	if err = db.Model(&orderShortShip).Updates(model.OrderShortShip{ShipperConfirmed: newName, Finished: true}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// CancelOrderShortShipHandler in database
func CancelOrderShortShipHandler(c *gin.Context) {
	orderShortShip, err := getOrderShortShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	canceleReason := c.PostForm("canceled_reason")

	orderShortShip.ID = getIDFromParam(c)
	if err = db.Model(&orderShortShip).Updates(model.OrderShortShip{Canceled: true, CanceledReason: canceleReason}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteOrderShortShipHandler in database
func DeleteOrderShortShipHandler(c *gin.Context) {
	if _, err := getOrderShortShipOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&model.OrderShortShip{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}

// -------------------- ORDER SHORT SHIP INTERNAL CALL FUNTION --------------------

// CreateOrderShortShip in database
func CreateOrderShortShip(orderID uint, shipperReceiveMoney bool) (uint, error) {

	orderShortShip := &model.OrderShortShip{}
	orderShortShip.OrderID = orderID
	orderInfoForShipment, err := getOrderInfoOrNotFoundForShipment(orderID)
	if err != nil {
		return uint(0), err
	}

	transportType := &model.TransportType{}
	if err := db.First(transportType, orderInfoForShipment.TransportTypeID).Error; err != nil {
		return uint(0), err
	}
	deliveryLocation := &model.DeliveryLocation{}

	if orderInfoForShipment.UseLongShip == true {
		if err := db.Where("city = ?", transportType.LocationTwo).First(deliveryLocation).Error; err != nil {
			return uint(0), err
		}
	} else {
		if err := db.Where("city = ?", transportType.LocationOne).First(deliveryLocation).Error; err != nil {
			return uint(0), err
		}
	}
	employeeList := []model.EmployeeInfoForShortShip{}
	selectPart := "e.id, e.employee_type_id, e.delivery_location_id "
	err = db.Table("employees as e").Select(selectPart).
		Where("e.employee_type_id = ? AND e.delivery_location_id", 2, deliveryLocation.ID).Find(employeeList).Error
	if err != nil {
		return uint(0), err
	}
	length := len(employeeList)
	index := rand.Intn(length - 1)
	orderShortShip.ShipperID = employeeList[index].ID

	orderShortShip.CustomerReceiveID = orderInfoForShipment.CustomerReceiveID
	orderShortShip.CustomerSendFCMToken = orderInfoForShipment.CustomerSendFCMToken
	orderShortShip.CustomerRecvFCMToken = orderInfoForShipment.CustomerRecvFCMToken
	orderShortShip.CustomerSendFCMToken = orderInfoForShipment.CustomerRecvFCMToken
	orderShortShip.ShipperReceiveMoney = shipperReceiveMoney
	if err := db.Create(orderShortShip).Error; err != nil {
		return uint(0), err
	}
	return orderShortShip.ID, nil
}