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

func createCustomerNotificationLongShipHandler(longShipID uint, titleIndex int, employeeID uint) error {
	var title string
	var content string
	var employyIDString = strconv.FormatUint(uint64(employeeID), 10)
	switch titleIndex {
	case 1:
		title = "Your package has loaded on truck"
		content = "Employee load ID: " + employyIDString
	case 2:
		title = "Your long ship truck has started"
		content = "Employee driver 1 ID: " + employyIDString
	case 3:
		title = "Your long ship truck has arrived"
		content = "Employee driver 2 ID: " + employyIDString
	case 4:
		title = "Your package has unloaded off truck"
		content = "Employee unload ID: " + employyIDString
	default:
		return errors.New("Server error")
	}
	customerNotification := &model.CustomerNotification{Title: title, Content: content}
	orderLongShips := []model.OrderLongShip{}
	if err := db.Where("long_ship_id == ?", longShipID).Find(&orderLongShips).Error; err != nil {
		return err
	}

	for i := 0; i < len(orderLongShips); i++ {
		customerNotification.CustomerID = orderLongShips[i].CustomerSendID
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
		title = "Selected shipper for your order"
		content = "Our shipper has received your order with id: " + orderIDString
	case 2:
		title = "Shipper has called you"
		content = "Your order id: " + orderIDString + " has been verified"
	case 3:
		title = "Shipper received your money"
		content = "Thanks for paying your order with id: " + orderIDString
	case 4:
		title = "Shipper shipped your package"
		content = "Your order id: " + orderIDString + " has arrived"
	case 5:
		title = "Shipper confirmed your package"
		content = "Thanks for using our service"
	case 6:
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
