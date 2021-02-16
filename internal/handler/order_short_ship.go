package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	CommonService "github.com/lucthienbinh/golang_scem/internal/service/common"
	CommonMessage "github.com/lucthienbinh/golang_scem/internal/service/common_message"
	"golang.org/x/sync/errgroup"
)

// -------------------- ORDER SHORT SHIP HANDLER FUNTION --------------------

// GetOrderShortShipListHandler in database
func GetOrderShortShipListHandler(c *gin.Context) {
	orderShortShips := []model.OrderShortShip{}
	db.Order("id asc").Find(&orderShortShips)
	c.JSON(http.StatusOK, gin.H{"order_short_ship_list": &orderShortShips})
	return
}

// GetOrderShortShipListByEmployeeIDHandler in database
func GetOrderShortShipListByEmployeeIDHandler(c *gin.Context) {

	type APIOrderShortShipList struct {
		ID        uint   `json:"id"`
		OrderID   uint   `json:"order_id"`
		ShipperID uint   `json:"shipper_id"`
		Sender    string `json:"sender"`
		Receiver  string `json:"receiver"`
		Canceled  bool   `json:"canceled"`
		Finished  bool   `json:"finished"`
		CreatedAt int64  `json:"created_at"`
	}
	orderShortShipList := []APIOrderShortShipList{}
	db.Model(&model.OrderShortShip{}).Order("id asc").Find(&orderShortShipList, "shipper_id = ?", c.Param("id"))

	c.JSON(http.StatusOK, gin.H{"order_short_ship_list": &orderShortShipList})
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
	c.JSON(http.StatusOK, gin.H{"order_short_ship_info": &orderShortShip})
	return
}

// UpdateOSSShipperCalledHandler in database
func UpdateOSSShipperCalledHandler(c *gin.Context) {
	orderShortShip, err := getOrderShortShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if orderShortShip.ShipperCalled == true || orderShortShip.Canceled == true {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Run concurrency
	var g errgroup.Group

	// Send message to state machine (Zeebe or State Scem)
	g.Go(func() error {
		if err := CommonMessage.PublishShipperCalledMessage(orderShortShip.OrderID); err != nil {
			return err
		}
		return nil
	})

	// Send notification to FCM cloud and store in database
	g.Go(func() error {
		if err := createCustomerNotificationLShortShipHandler(orderShortShip.OrderID, orderShortShip.CustomerSendID, 1, ""); err != nil {
			return err
		}
		return nil
	})

	// Update order long ship in database
	g.Go(func() error {
		if err := db.Model(&orderShortShip).Updates(model.OrderShortShip{ShipperCalled: true}).Error; err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
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

	if orderShortShip.ShipperReceiveMoney == false || orderShortShip.ShipperReceivedMoney == true || orderShortShip.ShipperCalled == false || orderShortShip.Canceled == true {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Run concurrency
	var g errgroup.Group

	// Send message to state machine (Zeebe or State Scem)
	g.Go(func() error {
		if err := CommonMessage.PublishShipperReceivedMoneyMessage(orderShortShip.OrderID); err != nil {
			return err
		}
		return nil
	})

	// Update order long ship in database
	g.Go(func() error {
		if err := db.Model(&orderShortShip).Updates(model.OrderShortShip{
			ShipperReceivedMoney: true,
			ReceivedMoneyTime:    time.Now().Unix(),
		}).Error; err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
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

	// Run concurrency
	var g errgroup.Group

	// Send message to state machine (Zeebe or State Scem)
	g.Go(func() error {
		if err := CommonMessage.PublishShipperShippedMessage(orderShortShip.OrderID); err != nil {
			return err
		}
		return nil
	})

	// Update order long ship in database
	g.Go(func() error {
		if err := db.Model(&orderShortShip).Updates(model.OrderShortShip{ShipperShipped: true, ShippedTime: time.Now().Unix()}).Error; err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
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
	newName := fmt.Sprintf("%x", b)
	createTime := fmt.Sprintf("%d", time.Now().Unix())
	newName = createTime + "_" + newName + ".jpg"
	filepath := os.Getenv("IMAGE_FILE_PATH") + newName

	// Upload the file to specific dst.
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Run concurrency
	var g errgroup.Group

	// Send message to state machine (Zeebe or State Scem)
	g.Go(func() error {
		if err := CommonMessage.PublishShipperConfirmedMessage(orderShortShip.OrderID); err != nil {
			return err
		}
		return nil
	})

	// Send notification to FCM cloud and store in database
	g.Go(func() error {
		if err := createCustomerNotificationLShortShipHandler(orderShortShip.OrderID, orderShortShip.CustomerSendID, 2, ""); err != nil {
			return err
		}
		return nil
	})

	// Update order long ship in database
	g.Go(func() error {
		if err := db.Model(&orderShortShip).Updates(model.OrderShortShip{
			ShipperConfirmed:     newName,
			ShipperConfirmedTime: time.Now().Unix(),
			Finished:             true,
		}).Error; err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
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

	orderWorkflowData, err := getOrderWorkflowDataByOrderID(orderShortShip.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	canceleReason := c.PostForm("canceled_reason")

	// Run concurrency
	var g errgroup.Group

	// Send message to state machine (Zeebe or State Scem)
	g.Go(func() error {
		if err := CommonService.CancelWorkflowFullShipInstanceHandler(orderWorkflowData.WorkflowInstanceKey); err != nil {
			return err
		}
		return nil
	})

	// Send notification to FCM cloud and store in database
	g.Go(func() error {
		if err := createCustomerNotificationLShortShipHandler(orderShortShip.OrderID, orderShortShip.CustomerSendID, 3, canceleReason); err != nil {
			return err
		}
		return nil
	})

	// Update order long ship in database
	g.Go(func() error {
		if err := db.Model(&orderShortShip).Updates(model.OrderShortShip{Canceled: true, CanceledReason: canceleReason}).Error; err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
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
