package handler

import (
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	SSWorkflow "github.com/lucthienbinh/golang_scem/internal/service/state_scem/workflow"
	ZBWorkflow "github.com/lucthienbinh/golang_scem/internal/service/zeebe/workflow"
)

// CreateWorkflowFullShipInstanceHandler will select private function
func CreateWorkflowFullShipInstanceHandler(orderWorkflowData *model.OrderWorkflowData) (string, uint, error) {
	if os.Getenv("STATE_SERVICE") == "1" {
		return createWorkflowFullShipInstanceHandlerZB(orderWorkflowData)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return createWorkflowFullShipInstanceHandlerSS(orderWorkflowData)
	}
	return "231-321314-41515135131", uint(1), nil
}

// CreateWorkflowInstanceSS function
func CreateWorkflowInstanceSS(c *gin.Context) {
	orderWorkflowData := model.OrderWorkflowData{}
	orderWorkflowData.OrderID = uint(12)
	orderWorkflowData.ShipperReceiveMoney = true
	orderWorkflowData.UseLongShip = true
	orderWorkflowData.CustomerReceiveID = uint(11)

	WorkflowKey, WorkflowInstanceKey, err := createWorkflowFullShipInstanceHandlerSS(&orderWorkflowData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response1": WorkflowKey, "server_response2": WorkflowInstanceKey})
	return
}

// ------------------------------ CALL TO ZEEBE CLIENT ------------------------------

// DeployWorkflowFullShipHandlerZB function
func DeployWorkflowFullShipHandlerZB(c *gin.Context) {
	ZBWorkflow.DeployFullShipWorkflow()
}

// DeployWorkflowLongShipHandlerZB function
func DeployWorkflowLongShipHandlerZB(c *gin.Context) {
	ZBWorkflow.DeployLongShipWorkflow()
}

func createWorkflowFullShipInstanceHandlerZB(orderWorkflowData *model.OrderWorkflowData) (string, uint, error) {
	WorkflowKey, WorkflowInstanceKey, err := ZBWorkflow.CreateFullShipInstance(orderWorkflowData)
	if err != nil {
		return "", uint(0), err
	}
	return WorkflowKey, WorkflowInstanceKey, nil
}

func createWorkflowLongShipInstanceHandlerZB(orderWorkflowData *model.OrderWorkflowData) (string, uint, error) {
	WorkflowKey, WorkflowInstanceKey, err := ZBWorkflow.CreateLongShipInstance(orderWorkflowData)
	if err != nil {
		return "", uint(0), err
	}
	return WorkflowKey, WorkflowInstanceKey, nil
}

// ------------------------------ CALL TO GOLANG STATE SCEM CLIENT ------------------------------

// DeployWorkflowFullShipHandlerSS function
func DeployWorkflowFullShipHandlerSS(c *gin.Context) {
	if err := SSWorkflow.DeployFullShipWorkflow(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "A workflow model has been created!"})
	return
}

// DeployWorkflowLongShipHandlerSS function
func DeployWorkflowLongShipHandlerSS(c *gin.Context) {
	if err := SSWorkflow.DeployLongShipWorkflow(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "A workflow model has been created!"})
	return
}

func createWorkflowFullShipInstanceHandlerSS(orderWorkflowData *model.OrderWorkflowData) (string, uint, error) {
	WorkflowKey, WorkflowInstanceKey, err := SSWorkflow.CreateFullShipInstance(orderWorkflowData)
	if err != nil {
		return "", uint(0), err
	}
	return WorkflowKey, WorkflowInstanceKey, nil
}

// ------------------------------ CALL FROM CLIENT ------------------------------

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
