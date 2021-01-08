package handler

import (
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	ZBWorkflow "github.com/lucthienbinh/golang_scem/internal/service/zeebe/workflow"
)

// ------------------------------ CALL TO ZEEBE CLIENT ------------------------------

// DeployWorkflowFullShipHandler function
func DeployWorkflowFullShipHandler(c *gin.Context) {
	ZBWorkflow.DeployLongShipWorkflow()
}

// DeployWorkflowLongShipHandler function
func DeployWorkflowLongShipHandler(c *gin.Context) {
	ZBWorkflow.DeployLongShipWorkflow()
}

func createWorkflowFullShipInstanceHandler(orderWorkflowData *model.OrderWorkflowData) (uint, uint, error) {
	WorkflowKey, WorkflowInstanceKey, err := ZBWorkflow.CreateFullShipInstance(orderWorkflowData)
	if err != nil {
		return uint(0), uint(0), err
	}
	return WorkflowKey, WorkflowInstanceKey, nil
}

func createWorkflowLongShipInstanceHandler(orderWorkflowData *model.OrderWorkflowData) (uint, uint, error) {
	WorkflowKey, WorkflowInstanceKey, err := ZBWorkflow.CreateLongShipInstance(orderWorkflowData)
	if err != nil {
		return uint(0), uint(0), err
	}
	return WorkflowKey, WorkflowInstanceKey, nil
}

// ------------------------------ CALL FROM ZEEBE CLIENT ------------------------------

// ++++++++++++++++++++ Credit Payment Worker ++++++++++++++++++++

// CreditPaymentService function
func CreditPaymentService(orderID uint) (bool, error) {
	// Get order info for payment
	orderInfoForPayment, err := getOrderInfoOrNotFoundForPayment(orderID)
	if err != nil {
		return false, nil
	}
	// Compare customer balance and TotalPrice to decide if customer can use credit or not
	useCredit, err := compareCustomerBalanceVsTotalPrice(orderInfoForPayment.CustomerSendID, orderInfoForPayment.TotalPrice)
	if err != nil {
		return false, nil
	}
	return useCredit, nil
}

// GetOrderLongShipList function
func GetOrderLongShipList(longShipID uint) ([]model.OrderLongShip, error) {
	orderLongShips := []model.OrderLongShip{}
	if err := db.Where("long_ship_id == ?", longShipID).Order("id asc").Find(&orderLongShips).Error; err != nil {
		return nil, err
	}
	return orderLongShips, nil
}

// ++++++++++++++++++++ Order Long Ship Worker ++++++++++++++++++++

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
func CreateOrderShortShip(orderID uint, shipperReceiveMoney bool) (uint, error) {

	orderShortShip := &model.OrderShortShip{}
	orderShortShip.OrderID = orderID
	orderInfoForShipment, err := getOrderInfoOrNotFoundForShipment(orderID)
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
		Where("e.employee_type_id = ? AND e.delivery_location_id", 2, deliveryLocation.ID).Find(employeeList).Error
	if err != nil {
		return uint(0), err
	}
	length := len(employeeList)
	index := rand.Intn(length - 1)
	orderShortShip.ShipperID = employeeList[index].ID

	orderShortShip.CustomerReceiveID = orderInfoForShipment.CustomerReceiveID
	orderShortShip.CustomerSendFCMToken = orderInfoForShipment.CustomerSendFCMToken
	orderShortShip.CustomerRecvFCMToken = orderInfoForShipment.CustomerRecvFCMToken
	orderShortShip.CustomerSendFCMToken = orderInfoForShipment.CustomerRecvFCMToken
	orderShortShip.ShipperReceiveMoney = shipperReceiveMoney
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
