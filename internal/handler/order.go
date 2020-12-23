package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
)

// -------------------- ORDER HANDLER FUNTION --------------------

// GetOrderInfoListHandler in database
func GetOrderInfoListHandler(c *gin.Context) {
	orderInfoList := []model.OrderInfoFetchDB{}
	selectPart := "ord.id, ord.weight, ord.volume, ord.type, ord.image, ord.has_package, " +
		"c1.name as customer_send_name, c2.name as customer_receive_name, t.name as trasnport_type, " +
		"e1.name as empl_create_name, e2.name as empl_ship_name, ord.receiver, ord.detail, ord.total_price, ord.note, ord.created_at"
	leftJoin1 := "left join customers as c1 on ord.customer_send_id = c1.id"
	leftJoin2 := "left join customers as c2 on ord.customer_receive_id = c2.id"
	leftJoin3 := "left join transport_types as t on ord.trasnport_type_id = t.id"
	leftJoin4 := "left join employees as e1 on ord.empl_create_id = e1.id"
	leftJoin5 := "left join employees as e2 on ord.empl_ship_id = e2.id"

	db.Table("order_infos as ord").Select(selectPart).Joins(leftJoin1).Joins(leftJoin2).Joins(leftJoin3).Joins(leftJoin4).Joins(leftJoin5).
		Order("ord.id asc").Find(&orderInfoList)

	c.JSON(http.StatusOK, gin.H{"order_info_list": orderInfoList})
	return
}

func getOrderInfoOrNotFound(c *gin.Context) (*model.OrderInfoFetchDB, error) {
	orderInfoFetchDB := &model.OrderInfoFetchDB{}
	selectPart := "ord.id, ord.weight, ord.volume, ord.type, ord.image, ord.has_package, " +
		"c1.name as customer_send_name, c2.name as customer_receive_name, t.name as trasnport_type, " +
		"e1.name as empl_create_name, e2.name as empl_ship_name, ord.receiver, ord.detail, ord.total_price, ord.note"
	leftJoin1 := "left join customers as c1 on ord.customer_send_id = c1.id"
	leftJoin2 := "left join customers as c2 on ord.customer_receive_id = c2.id"
	leftJoin3 := "left join transport_types as t on ord.trasnport_type_id = t.id"
	leftJoin4 := "left join employees as e1 on ord.empl_create_id = e1.id"
	leftJoin5 := "left join employees as e2 on ord.empl_ship_id = e2.id"

	if err := db.Table("order_infos as ord").Select(selectPart).Joins(leftJoin1).Joins(leftJoin2).Joins(leftJoin3).Joins(leftJoin4).Joins(leftJoin5).
		Order("ord.id asc").First(&orderInfoFetchDB, c.Param("id")).Error; err != nil {
		return orderInfoFetchDB, err
	}
	return orderInfoFetchDB, nil
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

// CreateOrderInfoHandler in database
func CreateOrderInfoHandler(c *gin.Context) {
	orderInfo := &model.OrderInfo{}
	if err := c.ShouldBindJSON(&orderInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&orderInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "An order info has been created!"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// // -------------------- TRANSPORT TYPE HANDLER FUNTION --------------------

// GetTransportTypeListHandler in database
func GetTransportTypeListHandler(c *gin.Context) {
	transportTypes := []model.TransportType{}
	db.Order("id asc").Find(&transportTypes)
	c.JSON(http.StatusOK, gin.H{"transport_type_list": &transportTypes})
	return
}

func getTransportTypeOrNotFound(c *gin.Context) (*model.TransportType, error) {
	transportType := &model.TransportType{}
	if err := db.First(&transportType, c.Param("id")).Error; err != nil {
		return transportType, err
	}
	return transportType, nil
}

// GetTransportTypeHandler in database
func GetTransportTypeHandler(c *gin.Context) {
	transportType, err := getTransportTypeOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transport_type_info": &transportType})
	return
}

// CreateTransportTypeHandler in database
func CreateTransportTypeHandler(c *gin.Context) {
	transportType := &model.TransportType{}
	if err := c.ShouldBindJSON(&transportType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&transportType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "An transport type has been created!"})
	return
}

// UpdateTransportTypeHandler in database
func UpdateTransportTypeHandler(c *gin.Context) {
	transportType, err := getTransportTypeOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&transportType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transportType.ID = getIDFromParam(c)
	if err = db.Model(&transportType).Updates(&transportType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteTransportTypeHandler in database
func DeleteTransportTypeHandler(c *gin.Context) {
	if _, err := getTransportTypeOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&model.TransportType{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}
