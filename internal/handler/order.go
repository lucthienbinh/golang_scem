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
	orderInfoList := []model.OrderInfoFetchDB{}
	selectPart := "ord.id, ord.weight, ord.volume, ord.type, ord.image, " +
		"c1.name as customer_send_name, c2.name as customer_receive_name, t.name as transport_type, " +
		"e1.name as empl_create_name, e2.name as empl_ship_name, ord.receiver, ord.detail, ord.total_price, ord.note, ord.created_at"
	leftJoin1 := "left join customers as c1 on ord.customer_send_id = c1.id"
	leftJoin2 := "left join customers as c2 on ord.customer_receive_id = c2.id"
	leftJoin3 := "left join transport_types as t on ord.transport_type_id = t.id"
	leftJoin4 := "left join employees as e1 on ord.empl_create_id = e1.id"
	leftJoin5 := "left join employees as e2 on ord.empl_ship_id = e2.id"

	db.Table("order_infos as ord").Select(selectPart).Joins(leftJoin1).Joins(leftJoin2).Joins(leftJoin3).Joins(leftJoin4).Joins(leftJoin5).
		Order("ord.id asc").Find(&orderInfoList)

	c.JSON(http.StatusOK, gin.H{"order_info_list": orderInfoList})
	return
}

func getOrderInfoOrNotFound(c *gin.Context) (*model.OrderInfoFetchDB, error) {
	orderInfoFetchDB := &model.OrderInfoFetchDB{}
	selectPart := "ord.id, ord.weight, ord.volume, ord.type, ord.image, " +
		"c1.name as customer_send_name, c2.name as customer_receive_name, t.name as transport_type, " +
		"e1.name as empl_create_name, e2.name as empl_ship_name, ord.receiver, ord.detail, ord.total_price, ord.note"
	leftJoin1 := "left join customers as c1 on ord.customer_send_id = c1.id"
	leftJoin2 := "left join customers as c2 on ord.customer_receive_id = c2.id"
	leftJoin3 := "left join transport_types as t on ord.transport_type_id = t.id"
	leftJoin4 := "left join employees as e1 on ord.empl_create_id = e1.id"
	leftJoin5 := "left join employees as e2 on ord.empl_ship_id = e2.id"

	if err := db.Table("order_infos as ord").Select(selectPart).Joins(leftJoin1).Joins(leftJoin2).Joins(leftJoin3).Joins(leftJoin4).Joins(leftJoin5).
		Order("ord.id asc").First(orderInfoFetchDB, c.Param("id")).Error; err != nil {
		return orderInfoFetchDB, err
	}
	return orderInfoFetchDB, nil
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
	orderInfoFetchDB, err := getOrderInfoOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order_info": &orderInfoFetchDB})
	return
}

func calculateTotalPrice(orderInfo *model.OrderInfo) (int64, error) {
	transportType := &model.TransportType{}
	if err := db.First(transportType, orderInfo.TransportTypeID).Error; err != nil {
		return 0, err
	}
	var totalPrice int64
	if orderInfo.UseLongShip == true {
		totalPrice += transportType.LongShipPrice
	}
	if orderInfo.UseShortShip == true {
		totalPrice += (transportType.ShortShipPricePerKm * orderInfo.ShortShipDistance)
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
