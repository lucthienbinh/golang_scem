package handler

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucthienbinh/golang_scem/internal/model"
	"golang.org/x/sync/errgroup"
	"gopkg.in/validator.v2"
)

// -------------------- ORDER HANDLER FUNTION --------------------

// GetOrderInfoListHandler in database
func GetOrderInfoListHandler(c *gin.Context) {

	type APIOrderList struct {
		ID                uint   `json:"id"`
		CustomerSendID    uint   `json:"customer_send_id" validate:"nonzero"`
		CustomerReceiveID uint   `json:"customer_receive_id"`
		Sender            string `json:"sender" validate:"nonzero"`
		Receiver          string `json:"receiver" validate:"nonzero"`
		TransportTypeID   uint   `json:"transport_type_id" validate:"nonzero"`
		TotalPrice        int64  `json:"total_price"`
	}
	orderInfoList := []APIOrderList{}
	db.Model(&model.OrderInfo{}).Order("id asc").Find(&orderInfoList)

	type APIOrderPay struct {
		ID                  uint   `gorm:"primary_key;<-:false" json:"id"`
		OrderID             uint   `json:"order_id"`
		PayMethod           string `json:"pay_method"`
		PayStatus           bool   `json:"pay_status"`
		TotalPrice          int64  `json:"total_price"`
		FinishedStepOne     bool   `json:"finished_step_one"`
		FinishedStepTwo     bool   `json:"finished_step_two"`
		ShipperReceiveMoney bool   `json:"shipper_receive_money"`
	}
	orderPays := []APIOrderPay{}
	db.Model(&model.OrderPay{}).Order("id asc").Find(&orderPays)

	c.JSON(http.StatusOK, gin.H{"order_info_list": orderInfoList, "order_pay_list": &orderPays})
	return
}

// GetOrderListByCustomerIDHandler in database
func GetOrderListByCustomerIDHandler(c *gin.Context) {

	type APIOrderList struct {
		ID             uint   `json:"id"`
		CustomerSendID uint   `json:"customer_send_id"`
		Receiver       string `json:"receiver" validate:"nonzero"`
		CreatedAt      int64  `json:"created_at"`
		Detail         string `json:"detail"`
		TotalPrice     int64  `json:"total_price"`
		Image          string `json:"image"`
	}
	orderInfoList := []APIOrderList{}
	db.Model(&model.OrderInfo{}).Order("id asc").Find(&orderInfoList, "customer_send_id = ?", c.Param("id"))

	c.JSON(http.StatusOK, gin.H{"order_info_list": orderInfoList})
	return
}

func getOrderInfoOrNotFound(c *gin.Context) (*model.OrderInfo, error) {
	orderInfo := &model.OrderInfo{}
	if err := db.First(orderInfo, c.Param("id")).Error; err != nil {
		return orderInfo, err
	}
	return orderInfo, nil
}

func getOrderInfoOrNotFoundForPayment(orderID uint) (*model.OrderInfoForPayment, error) {
	orderInfoForPayment := &model.OrderInfoForPayment{}
	err := db.Model(&model.OrderInfo{}).First(orderInfoForPayment, orderID).Error
	if err != nil {
		return orderInfoForPayment, err
	}
	return orderInfoForPayment, nil
}

func getOrderInfoOrNotFoundForShipment(orderID uint) (*model.OrderInfoForShipment, error) {
	orderInfoForShipment := &model.OrderInfoForShipment{}
	err := db.Model(&model.OrderInfo{}).Order("id asc").First(&orderInfoForShipment, orderID).Error
	if err != nil {
		return orderInfoForShipment, err
	}
	return orderInfoForShipment, nil
}

func getOrderWorkflowDataByOrderID(orderID uint) (*model.OrderWorkflowData, error) {

	orderWorkflowData := &model.OrderWorkflowData{}
	err := db.Model(orderWorkflowData).Order("id asc").First(orderWorkflowData, "order_id = ?", orderID).Error
	if err != nil {
		return orderWorkflowData, err
	}
	return orderWorkflowData, nil
}

// GetOrderInfoHandler in database
func GetOrderInfoHandler(c *gin.Context) {
	orderInfo, err := getOrderInfoOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order_info": &orderInfo})
	return
}

// CreateOrderFormData function
func CreateOrderFormData(c *gin.Context) {
	type APILongShipList struct {
		ID                       uint   `json:"id"`
		TransportTypeID          uint   `json:"transport_type_id"`
		LicensePlate             string `json:"license_plate"`
		EstimatedTimeOfDeparture int64  `json:"estimated_time_of_departure"`
		EstimatedTimeOfArrival   int64  `json:"estimated_time_of_arrival"`
		Finished                 bool   `json:"finished"`
	}
	longShips := []APILongShipList{}
	transportTypes := []model.TransportType{}
	// Run concurrency
	var g errgroup.Group

	// Get long ship list
	g.Go(func() error {
		if err := db.Model(&model.LongShip{}).Order("id asc").Find(&longShips, "finished is ? and estimated_time_of_departure > ?", false, time.Now().Unix()).Error; err != nil {
			return err
		}
		return nil
	})

	// Calculate total price and Create order ID base on Time
	g.Go(func() error {
		if err := db.Where("same_city is ?", false).Order("id asc").Find(&transportTypes).Error; err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"long_ship_list":      &longShips,
		"transport_type_list": &transportTypes,
	})
	return
}

// ImageOrderHandler updload image of employee
func ImageOrderHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Println(file.Filename)

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

	c.JSON(http.StatusCreated, gin.H{"filename": newName})
	return
}

func calculateTotalPrice(orderInfo *model.OrderInfo) (int64, error) {
	transportType := &model.TransportType{}
	if err := db.First(transportType, orderInfo.TransportTypeID).Error; err != nil {
		return 0, err
	}
	var totalPrice int64
	totalPrice = transportType.ShortShipPricePerKm * orderInfo.ShortShipDistance
	if orderInfo.UseLongShip == true {
		totalPrice += transportType.LongShipPrice
	}
	return totalPrice, nil
}

// CreateOrderInfoHandler in database
func CreateOrderInfoHandler(c *gin.Context) {
	orderInfo := &model.OrderInfo{}
	if err := c.ShouldBindJSON(&orderInfo); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&orderInfo); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Run concurrency
	var g errgroup.Group

	// Check customer sender id and customer receive id
	g.Go(func() error {
		if err := db.First(&model.Customer{}, orderInfo.CustomerSendID).Error; err != nil {
			return err
		}
		if orderInfo.CustomerReceiveID > 0 {
			if err := db.First(&model.Customer{}, orderInfo.CustomerReceiveID).Error; err != nil {
				return err
			}
		}
		return nil
	})

	// Calculate total price and Create order ID base on Time
	g.Go(func() error {
		var err error
		orderInfo.TotalPrice, err = calculateTotalPrice(orderInfo)
		if err != nil {
			return err
		}
		current := uuid.New().Time()
		currentString := fmt.Sprintf("%d", current)
		rawUint, _ := strconv.ParseUint(currentString, 10, 64)
		orderInfo.ID = uint(rawUint / 100000000000)
		return nil
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create order info
	if err := db.Create(orderInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"server_response": "An order info has been created!",
		"order_id":        orderInfo.ID,
		"total_price":     orderInfo.TotalPrice,
	})
	return
}

// CreateOrderUseVoucherInfoHandler in database
func CreateOrderUseVoucherInfoHandler(c *gin.Context) {

	orderInfoWithVoucher := &model.OrderInfoWithVoucher{}
	if err := c.ShouldBindJSON(&orderInfoWithVoucher); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&orderInfoWithVoucher); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	orderInfo, orderVoucherID := orderInfoWithVoucher.ConvertToBasicOrder()

	// Run concurrency
	var g errgroup.Group
	var voucherDiscount int64

	// Check customer sender id and customer receive id
	g.Go(func() error {
		if err := db.First(&model.Customer{}, orderInfo.CustomerSendID).Error; err != nil {
			return err
		}
		return nil
	})

	// Check voucher id and get voucher discount
	g.Go(func() error {
		if orderVoucherID != 0 {
			orderVoucher := &model.OrderVoucher{}
			if err := db.Find(&orderVoucher, "id = ? AND start_date < ? AND end_date > ?", orderVoucherID, time.Now().Unix(), time.Now().Unix()).Error; err != nil {
				return err
			}
			if *orderVoucher == (model.OrderVoucher{}) {
				return errors.New("Bad request sent to server")
			}
			voucherDiscount = orderVoucher.Discount
			return nil
		}
		voucherDiscount = 0
		return nil
	})

	// Calculate total price and Create order ID base on Time
	g.Go(func() error {
		var err error
		orderInfo.TotalPrice, err = calculateTotalPrice(orderInfo)
		if err != nil {
			return err
		}
		current := uuid.New().Time()
		currentString := fmt.Sprintf("%d", current)
		rawUint, _ := strconv.ParseUint(currentString, 10, 64)
		orderInfo.ID = uint(rawUint / 100000000000)
		return nil
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderInfo.TotalPrice -= voucherDiscount

	// Create order info
	if err := db.Create(orderInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"server_response": "An order info has been created!",
		"order_id":        orderInfo.ID,
		"total_price":     orderInfo.TotalPrice,
	})
	return
}

// UpdateOrderInfoHandler in database
func UpdateOrderInfoHandler(c *gin.Context) {
	orderInfo, err := getOrderInfoOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&orderInfo); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	orderInfo.ID = getIDFromParam(c)
	if err = db.Model(&orderInfo).Updates(&orderInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteOrderInfoHandler in database
func DeleteOrderInfoHandler(c *gin.Context) {
	if _, err := getOrderInfoOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&model.OrderInfo{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}
