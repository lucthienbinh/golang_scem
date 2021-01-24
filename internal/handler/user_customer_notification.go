package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
)

// -------------------- DELIVERY LOCATION HANDLER FUNTION --------------------

// GetCustomerNotificationListByCustomerIDHandler in database
func GetCustomerNotificationListByCustomerIDHandler(c *gin.Context) {
	customerNotification := []model.CustomerNotification{}
	db.Order("id desc").Find(&customerNotification, "customer_id = ?", c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"customer_notification_list": &customerNotification})
	return
}

func createCustomerNotificationLongShipHandler(longShipID uint, titleIndex int) error {
	var title string
	var orderIDString string
	var longIDString = strconv.FormatUint(uint64(longShipID), 10)

	switch titleIndex {
	case 1:
		title = "Your order has started long ship trip"
	case 2:
		title = "Your order has finished long ship trip"
	default:
		return errors.New("Server error")
	}
	customerNotification := &model.CustomerNotification{Title: title}
	orderLongShips := []model.OrderLongShip{}
	if err := db.Where("long_ship_id == ?", longShipID).Find(&orderLongShips).Error; err != nil {
		return err
	}

	for i := 0; i < len(orderLongShips); i++ {
		orderIDString = strconv.FormatUint(uint64(orderLongShips[i].OrderID), 10)
		customerNotification.CustomerID = orderLongShips[i].CustomerSendID
		customerNotification.Content = "Order id: " + orderIDString + " Long ship id: " + longIDString
		if err := db.Create(customerNotification).Error; err != nil {
			return err
		}
		// Todo: send this notification ton gorush to push to FCM after created
	}

	return nil
}

func createCustomerNotificationLShortShipHandler(orderID, customerSendID uint, titleIndex int, canceledReason string) error {
	var title string
	var content string
	var orderIDString = strconv.FormatUint(uint64(orderID), 10)
	switch titleIndex {
	case 1:
		title = "Shipper has called you"
		content = "Your order id: " + orderIDString + " has been verified"
	case 2:
		title = "Shipper has confirmed your package"
		content = "Thanks for using our service. Finished order id: " + orderIDString
	case 3:
		title = "Your order has been canceled"
		content = canceledReason
	default:
		return errors.New("Server error")
	}
	customerNotification := &model.CustomerNotification{
		CustomerID: customerSendID, Title: title, Content: content,
	}
	if err := db.Create(customerNotification).Error; err != nil {
		return err
	}
	return nil
}
