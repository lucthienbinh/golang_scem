package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
)

// -------------------- DELIVERY LOCATION HANDLER FUNTION --------------------

var (
	longShipMessage = []string{
		"Your package has loaded on truck!",
		"Your long ship truck has started!",
		"Your long ship truck has arrived!",
		"Your package has unloaded off truck!",
	}
)

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
	}

	return nil
}

func createCustomerNotificationLShortShipHandler(customerID uint, titleIndex int, content string) error {
	var title string
	switch titleIndex {
	case 1:
		title = "Selected shipper for your order"
	case 2:
		title = "Shipper to verified"
	case 3:
		title = "Shipper received your money"
	case 4:
		title = "Shipper shipped your package"
	case 5:
		title = "Shipper confirmed your package"
	case 6:
		title = "Your order canceled"
	default:
		return errors.New("Server error")
	}
	customerNotification := &model.CustomerNotification{
		CustomerID: customerID, Title: title, Content: content,
	}
	if err := db.Create(customerNotification).Error; err != nil {
		return err
	}
	return nil
}
