package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
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
	c.JSON(http.StatusOK, gin.H{"order_short_ship_list": &orderShortShips})
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

// UpdateOSSShipperCalledHandler in database
func UpdateOSSShipperCalledHandler(c *gin.Context) {
	orderShortShip, err := getOrderShortShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if orderShortShip.ShipperCalled == true {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	orderShortShip.ID = getIDFromParam(c)
	if err = db.Model(&orderShortShip).Updates(model.OrderShortShip{ShipperCalled: true}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateOSSShipperReceivedMoneyHandler in database
func UpdateOSSShipperReceivedMoneyHandler(c *gin.Context) {
	orderShortShip, err := getOrderShortShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if orderShortShip.ShipperReceivedMoney == true || orderShortShip.ShipperCalled == false || orderShortShip.Canceled == true {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	orderShortShip.ID = getIDFromParam(c)
	if err = db.Model(&orderShortShip).Updates(model.OrderShortShip{
		ShipperReceivedMoney: true,
		ReceivedMoneyTime:    time.Now().Unix(),
	}).Error; err != nil {
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

	if orderShortShip.ShipperCalled == false || orderShortShip.ShipperShipped == true || orderShortShip.Canceled == true {
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

// UpdateOSSShipperConfirmedHandler in database
func UpdateOSSShipperConfirmedHandler(c *gin.Context) {
	orderShortShip, err := getOrderShortShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if orderShortShip.ShipperShipped == false || orderShortShip.ShipperConfirmed != "" || orderShortShip.Canceled == true {
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
	filepath := os.Getenv("IMAGE_FILE_PATH") + newName

	// Upload the file to specific dst.
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderShortShip.ID = getIDFromParam(c)
	if err = db.Model(&orderShortShip).Updates(model.OrderShortShip{
		ShipperConfirmed:     newName,
		ShipperConfirmedTime: time.Now().Unix(),
		Finished:             true,
	}).Error; err != nil {
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

	if orderShortShip.Finished == true || orderShortShip.Canceled == true {
		c.AbortWithStatus(http.StatusBadRequest)
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
