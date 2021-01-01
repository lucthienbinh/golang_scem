package handler

import (
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
