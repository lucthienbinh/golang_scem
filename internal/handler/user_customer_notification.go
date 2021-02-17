package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	GorushClient "github.com/lucthienbinh/golang_scem/internal/service/gorush"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

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

	gorushStatus := os.Getenv("GORUSH_STATUS")
	if gorushStatus == "1" {
		for i := 0; i < len(orderLongShips); i++ {
			orderIDString = strconv.FormatUint(uint64(orderLongShips[i].OrderID), 10)
			customerNotification.CustomerID = orderLongShips[i].CustomerSendID
			customerNotification.Content = "Order id: " + orderIDString + " Long ship id: " + longIDString
			userFCMToken := &model.UserFCMToken{}
			loggedIn := false

			// Run concurrency
			var g errgroup.Group

			// Save to database
			g.Go(func() error {
				if err := db.Create(customerNotification).Error; err != nil {
					return err
				}
				return nil
			})

			// Send to gorush -> FCM cloud -> Android
			g.Go(func() error {
				err := db.Model(userFCMToken).Order("id asc").First(userFCMToken, "customer_id = ?", customerNotification.CustomerID).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil
				}
				loggedIn = true
				return err
			})

			if err := g.Wait(); err != nil {
				return err
			}
			if loggedIn == true {
				if err := GorushClient.Client(userFCMToken.Token, title, customerNotification.Content); err != nil {
					return err
				}
			}
		}
	} else {
		for i := 0; i < len(orderLongShips); i++ {
			orderIDString = strconv.FormatUint(uint64(orderLongShips[i].OrderID), 10)
			customerNotification.CustomerID = orderLongShips[i].CustomerSendID
			customerNotification.Content = "Order id: " + orderIDString + " Long ship id: " + longIDString
			if err := db.Create(customerNotification).Error; err != nil {
				return err
			}
		}
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

	gorushStatus := os.Getenv("GORUSH_STATUS")
	if gorushStatus == "1" {
		userFCMToken := &model.UserFCMToken{}
		// Run concurrency
		var g errgroup.Group
		loggedIn := false

		// Save to database
		g.Go(func() error {
			if err := db.Create(customerNotification).Error; err != nil {
				return err
			}
			return nil
		})

		// Send to gorush -> FCM cloud -> Android
		g.Go(func() error {
			err := db.Model(userFCMToken).Order("id asc").First(userFCMToken, "customer_id = ?", customerSendID).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			loggedIn = true
			return err
		})

		if err := g.Wait(); err != nil {
			return err
		}
		if loggedIn == true {
			if err := GorushClient.Client(userFCMToken.Token, title, content); err != nil {
				return err
			}
		}
	} else {
		if err := db.Create(customerNotification).Error; err != nil {
			return err
		}
	}
	return nil
}
