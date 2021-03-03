package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	CommonService "github.com/lucthienbinh/golang_scem/internal/service/common"
	CommonMessage "github.com/lucthienbinh/golang_scem/internal/service/common_message"
	"golang.org/x/sync/errgroup"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// -------------------- ORDER PAYMENT STRUCT FOR BINDING JSON --------------------

type createOrderPay struct {
	OrderID             uint   `json:"order_id" validate:"nonzero"`
	PayMethod           string `json:"pay_method"`
	ShipperReceiveMoney bool   `json:"shipper_receive_money"`
}

// -------------------- ORDER PAYMENT HANDLER FUNTION --------------------

// GetOrderPayListHandler in database
func GetOrderPayListHandler(c *gin.Context) {
	orderPays := []model.OrderPay{}
	db.Order("id asc").Find(&orderPays)
	c.JSON(http.StatusOK, gin.H{"order_pay_list": &orderPays})
	return
}

func getOrderPayOrNotFound(orderPayID uint) (*model.OrderPay, error) {
	orderPay := &model.OrderPay{}
	if err := db.First(orderPay, orderPayID).Error; err != nil {
		return orderPay, err
	}
	return orderPay, nil
}

func getOrderPayOrNotFoundByOrderID(orderID uint) (*model.OrderPay, error) {
	orderPay := &model.OrderPay{}
	if err := db.Where("order_id = ?", orderID).First(orderPay).Error; err != nil {
		return orderPay, err
	}
	return orderPay, nil
}

func compareCustomerBalanceVsTotalPrice(customerID uint, totalPrice int64) (bool, error) {
	customerCredit := &model.CustomerCredit{}
	if err := db.Where("customer_id = ?", customerID).First(customerCredit).Error; err != nil {
		return false, err
	}
	if customerCredit.AccountBalance > totalPrice {
		return true, nil
	}
	return false, nil
}

func updateCustomerCreditAfterConfirmed(customerID uint, totalPrice int64) error {
	customerCredit := &model.CustomerCredit{}
	if err := db.Where("customer_id = ?", customerID).First(customerCredit).Error; err != nil {
		return err
	}
	if customerCredit.AccountBalance < totalPrice {
		return errors.New("Not enough money")
	}
	customerNewBalance := customerCredit.AccountBalance - totalPrice
	if err := updateCustomerCreditBalance(customerNewBalance, customerID); err != nil {
		return err
	}
	return nil
}

// CreateOrderPayStepOneHandler in database
func CreateOrderPayStepOneHandler(c *gin.Context) {
	stepOneRequest := createOrderPay{}
	if err := c.ShouldBindJSON(&stepOneRequest); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&stepOneRequest); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Get order info for payment
	orderInfoForPayment, err := getOrderInfoOrNotFoundForPayment(stepOneRequest.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Run concurrency
	var g errgroup.Group
	var useCredit bool

	// Compare customer balance and TotalPrice to decide if customer can use credit or not
	g.Go(func() error {
		var err error
		useCredit, err = compareCustomerBalanceVsTotalPrice(orderInfoForPayment.CustomerSendID, orderInfoForPayment.TotalPrice)
		if err != nil {
			return err
		}
		return nil
	})

	// Check if order pay is created or not
	g.Go(func() error {
		orderPay, err := getOrderPayOrNotFoundByOrderID(stepOneRequest.OrderID)
		if err != nil {
			// if not we will create a new data
			if errors.Is(err, gorm.ErrRecordNotFound) {
				orderPay = &model.OrderPay{OrderID: stepOneRequest.OrderID, TotalPrice: orderInfoForPayment.TotalPrice, FinishedStepOne: true}
				if err := db.Create(orderPay).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return hideCreditButton = true if customer balance is lower than TotalPrice
	if useCredit == true {
		c.JSON(http.StatusCreated, gin.H{"hideCreditButton": false})
	} else {
		c.JSON(http.StatusCreated, gin.H{"hideCreditButton": true})
	}

	return
}

// CreateOrderPayStepTwoHandler in database
func CreateOrderPayStepTwoHandler(c *gin.Context) {
	stepTwoRequest := createOrderPay{}
	if err := c.ShouldBindJSON(&stepTwoRequest); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&stepTwoRequest); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Run concurrency
	var g errgroup.Group
	orderPay := &model.OrderPay{}
	orderInfoForPayment := &model.OrderInfoForPayment{}

	// Get order pay info for validate step 1
	g.Go(func() error {
		var err error
		orderPay, err = getOrderPayOrNotFoundByOrderID(stepTwoRequest.OrderID)
		if err != nil {
			return err
		}
		if orderPay.FinishedStepOne == false {
			return errors.New("Bad request sent to server")
		}
		return nil
	})

	// Get order info for payment
	g.Go(func() error {
		var err error
		orderInfoForPayment, err = getOrderInfoOrNotFoundForPayment(stepTwoRequest.OrderID)
		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// shipper_receive_money will be sent when using app
	orderWorkflowData := &model.OrderWorkflowData{
		OrderID:             stepTwoRequest.OrderID,
		ShipperReceiveMoney: stepTwoRequest.ShipperReceiveMoney,
		UseLongShip:         orderInfoForPayment.UseLongShip,
		CustomerSendID:      orderInfoForPayment.CustomerSendID,
		CustomerReceiveID:   orderInfoForPayment.CustomerReceiveID,
	}
	// Create workflow instance in zeebe
	WorkflowKey, WorkflowInstanceKey, err := CommonService.CreateWorkflowFullShipInstanceHandler(orderWorkflowData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	orderWorkflowData.OrderPayID = orderPay.ID
	orderWorkflowData.LongShipID = orderInfoForPayment.LongShipID
	orderWorkflowData.WorkflowKey = WorkflowKey
	orderWorkflowData.WorkflowInstanceKey = WorkflowInstanceKey
	// Create workflow data in database
	if err := db.Create(orderWorkflowData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update order pay step 2
	orderPay.PayMethod = stepTwoRequest.PayMethod
	orderPay.FinishedStepTwo = true
	if err := db.Model(&orderPay).Updates(&orderPay).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "New workflow instance has been created!"})
	return
}

// UpdateOrderPayConfirmHandler in database
func UpdateOrderPayConfirmHandler(c *gin.Context) {
	var orderID = getIDFromParam(c)
	// Run concurrency
	var g errgroup.Group
	orderPay := &model.OrderPay{}
	orderInfoForPayment := &model.OrderInfoForPayment{}

	// Update customer credit
	g.Go(func() error {
		var err error
		orderPay, err = getOrderPayOrNotFoundByOrderID(orderID)
		if err != nil {
			return err
		}
		if orderPay.FinishedStepTwo == false {
			c.AbortWithStatus(http.StatusBadRequest)
			return errors.New("Bad request sent to server")
		}
		orderInfoForPayment, err = getOrderInfoOrNotFoundForPayment(orderID)
		if err != nil {
			return err
		}
		orderPay.PayStatus = true
		if orderPay.PayMethod == "credit" {
			err = updateCustomerCreditAfterConfirmed(orderInfoForPayment.CustomerSendID, orderPay.TotalPrice)
			if err != nil {
				return err
			}
		}
		if err = db.Model(&orderPay).Updates(&orderPay).Error; err != nil {
			return err
		}
		return nil
	})

	// Update customer credit
	g.Go(func() error {
		err := CommonMessage.PublishPaymentConfirmedMessage(orderID)
		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"server_response": "Order payment money has been received and confirmed!"})
	return
}

// Todo: get user token first then go to handler
// UpdateOrderPayCustomerConfirmHandler in database
// func UpdateOrderPayCustomerConfirmHandler(c *gin.Context, userAuthID uint) {
// 	userAuthenticate := &model.UserAuthenticate{}
// 	if err := db.First(userAuthenticate, userAuthID).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
// 		return
// 	}
// 	orderPay, err := getOrderPayOrNotFound(getIDFromParam(c))
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 		return
// 	}
// 	orderPay.PayCustomerID = userAuthenticate.CustomerID
// 	if err := db.Model(&orderPay).Updates(&orderPay).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, gin.H{
// 		"server_response": "Order payment money has been confirmed by customer id: ",
// 		"employee_id":     userAuthenticate.CustomerID})
// 	return
// }
