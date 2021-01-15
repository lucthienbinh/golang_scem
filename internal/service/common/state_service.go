package common

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	SSWorkflow "github.com/lucthienbinh/golang_scem/internal/service/state_scem/workflow"
	ZBWorkflow "github.com/lucthienbinh/golang_scem/internal/service/zeebe/workflow"
)

///////////////////////////////////////// CREATE INSTANCE /////////////////////////////////////////

////////+++++++++++ SERVICE SELECTOR +++++++++++////////

// CreateWorkflowFullShipInstanceHandler will select private function
func CreateWorkflowFullShipInstanceHandler(orderWorkflowData *model.OrderWorkflowData) (string, uint, error) {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBWorkflow.CreateFullShipInstance(orderWorkflowData)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return SSWorkflow.CreateFullShipInstance(orderWorkflowData)
	}
	return "231-321314-41515135131", uint(1), nil
}

// Todo: Create full ship instance with state scem

// CreateWorkflowLongShipInstanceHandler will select private function
func CreateWorkflowLongShipInstanceHandler(longShipID uint) (string, uint, error) {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBWorkflow.CreateLongShipInstance(longShipID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return "231-321314-41515135131", uint(1), nil
	}
	return "231-321314-41515135131", uint(1), nil
}

///////////////////////////////////////// DEPLOY WORKFLOW /////////////////////////////////////////

// DeployWorkflowFullShipHandlerZB function
func DeployWorkflowFullShipHandlerZB(c *gin.Context) {
	ZBWorkflow.DeployFullShipWorkflow()
}

// DeployWorkflowLongShipHandlerZB function
func DeployWorkflowLongShipHandlerZB(c *gin.Context) {
	ZBWorkflow.DeployLongShipWorkflow()
}

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
