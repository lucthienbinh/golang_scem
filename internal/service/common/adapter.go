package common

import (
	"math/rand"

	"github.com/lucthienbinh/golang_scem/internal/model"
	"gorm.io/gorm"
)

var db *gorm.DB

// MappingGormDBConnection to open connect with database
func MappingGormDBConnection(_db *gorm.DB) {
	db = _db
}

// ------------------------------ CALL FROM STATE SERVICE ------------------------------

// GetOrderLongShipList function
func GetOrderLongShipList(longShipID uint) ([]model.OrderLongShip, error) {
	orderLongShips := []model.OrderLongShip{}
	if err := db.Where("long_ship_id == ?", longShipID).Order("id asc").Find(&orderLongShips).Error; err != nil {
		return nil, err
	}
	return orderLongShips, nil
}

// ++++++++++++++++++++ Order Long Ship Worker ++++++++++++++++++++

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

// CreateOrderLongShip in database
func CreateOrderLongShip(orderID uint) (uint, error) {

	orderInfoForShipment, err := getOrderInfoOrNotFoundForShipment(orderID)
	if err != nil {
		return uint(0), err
	}

	orderLongShip := &model.OrderLongShip{}
	orderLongShip.OrderID = orderID
	orderLongShip.LongShipID = orderInfoForShipment.LongShipID
	orderLongShip.CustomerSendFCMToken = orderInfoForShipment.CustomerSendFCMToken
	orderLongShip.CustomerRecvFCMToken = orderInfoForShipment.CustomerRecvFCMToken

	if err := db.Create(orderLongShip).Error; err != nil {
		return uint(0), err
	}

	orderInfo := &model.OrderInfo{}
	orderInfo.ID = orderID
	orderInfo.OrderLongShipID = orderLongShip.ID
	if err = db.Model(&orderInfo).Updates(&orderInfo).Error; err != nil {
		return uint(0), err
	}
	return orderLongShip.ID, nil
}

// ++++++++++++++++++++ Order Short Ship Worker ++++++++++++++++++++

// CreateOrderShortShip in database
func CreateOrderShortShip(orderID uint) (uint, error) {

	orderShortShip := &model.OrderShortShip{}
	orderShortShip.OrderID = orderID
	orderInfoForShipment, err := getOrderInfoOrNotFoundForShipment(orderID)
	if err != nil {
		return uint(0), err
	}
	orderWorkflowData, err := getOrderWorkflowDataByOrderID(orderID)
	if err != nil {
		return uint(0), err
	}

	transportType := &model.TransportType{}
	if err := db.First(transportType, orderInfoForShipment.TransportTypeID).Error; err != nil {
		return uint(0), err
	}
	deliveryLocation := &model.DeliveryLocation{}

	if orderInfoForShipment.UseLongShip == true {
		if err := db.Where("city = ?", transportType.LocationTwo).First(deliveryLocation).Error; err != nil {
			return uint(0), err
		}
	} else {
		if err := db.Where("city = ?", transportType.LocationOne).First(deliveryLocation).Error; err != nil {
			return uint(0), err
		}
	}
	// Pick random employee type driver base on location
	employeeList := []model.EmployeeInfoForShortShip{}
	selectPart := "e.id, e.employee_type_id, e.delivery_location_id "
	err = db.Table("employees as e").Select(selectPart).
		Where("e.employee_type_id = ? AND e.delivery_location_id", 2, deliveryLocation.ID).Find(&employeeList).Error
	if err != nil {
		return uint(0), err
	}
	length := len(employeeList)
	index := rand.Intn(length - 1)
	orderShortShip.ShipperID = employeeList[index].ID
	orderShortShip.CustomerSendID = orderInfoForShipment.CustomerSendID
	orderShortShip.CustomerReceiveID = orderInfoForShipment.CustomerReceiveID
	orderShortShip.CustomerSendFCMToken = orderInfoForShipment.CustomerSendFCMToken
	orderShortShip.CustomerRecvFCMToken = orderInfoForShipment.CustomerRecvFCMToken
	orderShortShip.ShipperReceiveMoney = orderWorkflowData.ShipperReceiveMoney
	if err := db.Create(orderShortShip).Error; err != nil {
		return uint(0), err
	}
	orderInfo := &model.OrderInfo{}
	orderInfo.ID = orderID
	orderInfo.OrderShortShipID = orderShortShip.ID
	if err = db.Model(&orderInfo).Updates(&orderInfo).Error; err != nil {
		return uint(0), err
	}
	return orderShortShip.ID, nil
}
