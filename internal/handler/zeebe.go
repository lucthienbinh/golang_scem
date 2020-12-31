package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	ZBWorkflow "github.com/lucthienbinh/golang_scem/internal/service/zeebe/workflow"
)

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
