package handler

import (
	"github.com/lucthienbinh/golang_scem/internal/model"
	ZBWorkflow "github.com/lucthienbinh/golang_scem/internal/service/zeebe/workflow"
)

func createWorkflowInstanceHandler(orderWorkflowData *model.OrderWorkflowData) (uint, uint, error) {
	WorkflowKey, WorkflowInstanceKey, err := ZBWorkflow.CreateNewInstance(orderWorkflowData)
	if err != nil {
		return uint(0), uint(0), err
	}
	return WorkflowKey, WorkflowInstanceKey, nil
}
