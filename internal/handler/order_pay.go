package handler

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// -------------------- ORDER PAYMENT STRUCT FOR BINDING JSON --------------------

type createOrderPay struct {
	OrderID             uint   `json:"order_id" validate:"nonzero"`
	PayMethod           string `json:"pay_method"`
	ShipperReceiveMoney bool   `json:"shipper_receive_money"`
}

type confirmOrderPay struct {
	ConfirmString string `json:"confirm_string" validate:"nonzero"`
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

// GetOrderPayHandler in database
func GetOrderPayHandler(c *gin.Context) {
	orderPay, err := getOrderPayOrNotFound(getIDFromParam(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order_pay_info": &orderPay})
	return
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
	// Compare customer balance and TotalPrice to decide if customer can use credit or not
	useCredit, err := compareCustomerBalanceVsTotalPrice(orderInfoForPayment.CustomerSendID, orderInfoForPayment.TotalPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Check if order pay is created or not
	orderPay, err := getOrderPayOrNotFoundByOrderID(stepOneRequest.OrderID)
	if err != nil {
		// if not we will create a new data
		if errors.Is(err, gorm.ErrRecordNotFound) {
			orderPay = &model.OrderPay{OrderID: stepOneRequest.OrderID, TotalPrice: orderInfoForPayment.TotalPrice, FinishedStepOne: true}
			if err := db.Create(orderPay).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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
	// Get order info for payment
	orderInfoForPayment, err := getOrderInfoOrNotFoundForPayment(stepTwoRequest.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	// Get order pay info for validate step 1 and create workflow instance
	orderPay, err := getOrderPayOrNotFoundByOrderID(stepTwoRequest.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if orderPay.FinishedStepOne == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	orderWorkflowData := &model.OrderWorkflowData{
		OrderID:             orderPay.OrderID,
		PayMethod:           stepTwoRequest.PayMethod,
		ShipperReceiveMoney: stepTwoRequest.ShipperReceiveMoney,
		UseLongShip:         orderInfoForPayment.UseLongShip,
		CustomerSendID:      orderInfoForPayment.CustomerSendID,
		CustomerReceiveID:   orderInfoForPayment.CustomerReceiveID,
	}
	// Create workflow instance in zeebe
	WorkflowKey, WorkflowInstanceKey, err := CreateWorkflowFullShipInstanceHandler(orderWorkflowData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	orderPay.ConfirmString = fmt.Sprintf("%x", b)
	if err := db.Model(&orderPay).Updates(&orderPay).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "New workflow instance has been created!"})
	return
}

// UpdateOrderPayEmployeeConfirmHandler in database
func UpdateOrderPayEmployeeConfirmHandler(c *gin.Context, userAuthID uint) {
	confirmString := c.PostForm("confirm_string")
	userAuthenticate := &model.UserAuthenticate{}
	if err := db.First(userAuthenticate, userAuthID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	orderPay, err := getOrderPayOrNotFound(getIDFromParam(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if confirmString != orderPay.ConfirmString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Confirm string needs to be match with the given!"})
		return
	}
	orderPay.PayEmployeeID = userAuthenticate.EmployeeID
	if err := db.Model(&orderPay).Updates(&orderPay).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"server_response": "Order payment money has been received and confirmed by employee id: ",
		"employee_id":     userAuthenticate.EmployeeID})
	return
}

// UpdateOrderPayCustomerConfirmHandler in database
func UpdateOrderPayCustomerConfirmHandler(c *gin.Context, userAuthID uint) {
	confirmString := c.PostForm("confirm_string")
	userAuthenticate := &model.UserAuthenticate{}
	if err := db.First(userAuthenticate, userAuthID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	orderPay, err := getOrderPayOrNotFound(getIDFromParam(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if confirmString != orderPay.ConfirmString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Confirm string needs to be match with the given!"})
		return
	}
	orderPay.PayCustomerID = userAuthenticate.CustomerID
	if err := db.Model(&orderPay).Updates(&orderPay).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"server_response": "Order payment money has been confirmed by customer id: ",
		"employee_id":     userAuthenticate.CustomerID})
	return
}
