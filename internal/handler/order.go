package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
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
	selectPart := "ord.id, ord.customer_send_id, ord.customer_receive_id, ord.use_long_ship, ord.use_short_ship, ord.total_price "
	err := db.Table("order_infos as ord").Select(selectPart).First(orderInfoForPayment, orderID).Error
	if err != nil {
		return orderInfoForPayment, err
	}
	return orderInfoForPayment, nil
}

func getOrderInfoOrNotFoundForShipment(orderID uint) (*model.OrderInfoForShipment, error) {
	orderInfoForShipment := &model.OrderInfoForShipment{}
	selectPart := "ord.id, ord.customer_send_fcm_token, ord.customer_recv_fcm_token, ord.long_ship_id, ord.customer_receive_id "
	err := db.Table("order_infos as ord").Select(selectPart).First(orderInfoForShipment, orderID).Error
	if err != nil {
		return orderInfoForShipment, err
	}
	return orderInfoForShipment, nil
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
	db.Model(&model.LongShip{}).Order("id asc").Find(&longShips)

	transportTypes := []model.TransportType{}
	db.Where("same_city is ?", false).Order("id asc").Find(&transportTypes)

	c.JSON(http.StatusOK, gin.H{
		"long_ship_list":      &longShips,
		"transport_type_list": &transportTypes,
	})
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
		// c.AbortWithStatus(http.StatusBadRequest)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate(&orderInfo); err != nil {
		// c.AbortWithStatus(http.StatusBadRequest)
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	// Calculate total price
	totalPrice, err := calculateTotalPrice(orderInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	orderInfo.TotalPrice = totalPrice
	// Insert customer FCM token to orderInfo to send message
	cusSendFCMToken := &model.UserFCMToken{}
	if err := db.Where("customer_id = ?", orderInfo.CustomerSendID).First(cusSendFCMToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cusReceiveFCMToken := &model.UserFCMToken{}
	if orderInfo.CustomerReceiveID != 0 {
		if err := db.Where("customer_id = ?", orderInfo.CustomerReceiveID).First(cusReceiveFCMToken).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orderInfo.CustomerRecvFCMToken = cusReceiveFCMToken.Token
	}
	orderInfo.CustomerSendFCMToken = cusSendFCMToken.Token
	// Create order info
	if err := db.Create(orderInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"server_response": "An order info has been created!",
		"order_id":        orderInfo.ID,
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
