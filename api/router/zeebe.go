package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/handler"
	"github.com/lucthienbinh/golang_scem/internal/model"
	ZBWorkflow "github.com/lucthienbinh/golang_scem/internal/service/zeebe/workflow"
)

// ZeebeRoutes to manage user model
func ZeebeRoutes(rg *gin.RouterGroup) {

	zeebe := rg.Group("/zeebe")
	zeebe.GET("/deploy-workflow", deployWorkflowHandler)
	zeebe.GET("/create-instance", createWorkflowInstanceHandler)
	zeebe.GET("/id/:id", handler.GetCustomerHandler)
	zeebe.POST("/create", handler.CreateCustomerHandler)
	zeebe.PUT("/update/:id", handler.UpdateCustomerHandler)
	zeebe.DELETE("/delete/:id", handler.DeleteCustomerHandler)

	orderLongShip := zeebe.Group("/order-long-ship")
	orderLongShip.GET("/list", handler.GetOrderLongShipListHandler)
	orderLongShip.GET("/id/:id", handler.GetOrderLongShipHandler)
	orderLongShip.PUT("/update/load-package", handler.UpdateOLSLoadPackageHandler)
	orderLongShip.PUT("/update/start-vehicle", handler.UpdateOLSStartVehicleHandler)
	orderLongShip.PUT("/update/vehicle-arrived", handler.UpdateOLSVehicleArrivedHandler)
	orderLongShip.PUT("/update/unload-package", handler.UpdateOLSUnloadPackageHandler)
	orderLongShip.DELETE("/delete/:id", handler.DeleteOrderLongShipHandler)

	orderShortShip := zeebe.Group("/order-short-ship")
	orderShortShip.GET("/list", handler.GetOrderShortShipListHandler)
	orderShortShip.GET("/id/:id", handler.GetOrderShortShipHandler)
	orderShortShip.PUT("/update/shipper-received", handler.UpdateOSSShipperReceivedHandler)
	orderShortShip.PUT("/update/shipper-called", handler.UpdateOSSShipperCalledHandler)
	orderShortShip.PUT("/update/shipper-shipped", handler.UpdateOSSShipperShippedHandler)
	orderShortShip.PUT("/update/cus-receive-confirmed", handler.UpdateOSSCusReceiveConfirmedHandler)
	orderShortShip.PUT("/update/shipper-confirmed", handler.UpdateOSSShipperConfirmedHandler)
	orderShortShip.PUT("/update/cancel-order", handler.CancelOrderShortShipHandler)
	orderShortShip.DELETE("/delete/:id", handler.DeleteOrderLongShipHandler)
}

func createWorkflowInstanceHandler(c *gin.Context) {
	orderWorkflowCreate := &model.OrderWorkflowCreate{}
	// if err := c.ShouldBindJSON(&orderWorkflowCreate); err != nil {
	// 	c.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }
	// if err := validator.Validate(&orderWorkflowCreate); err != nil {
	// 	c.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }
	orderWorkflowCreate.OrderID = 2
	orderWorkflowCreate.PayMethod = "credit"
	orderWorkflowCreate.TotalPrice = 12000
	orderWorkflowCreate.UseLongShip = true
	orderWorkflowCreate.UseShortShip = true
	orderPayID, err := handler.CreateOrderPayHandler(orderWorkflowCreate.OrderID, orderWorkflowCreate.PayMethod, orderWorkflowCreate.TotalPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ZBWorkflow.CreateNewInstance(1, 2, "credit", true, true)
	return
}

func deployWorkflowHandler(c *gin.Context) {
	ZBWorkflow.DeployNewWorkflow()
}
