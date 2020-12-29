package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/api/middleware"
	"github.com/lucthienbinh/golang_scem/internal/handler"
)

// OrderShipRoutes to manage user model
func OrderShipRoutes(rg *gin.RouterGroup) {

	zeebe := rg.Group("/zeebe")
	zeebe.GET("/deploy-workflow", handler.DeployWorkflowHandler)

	orderLongShip := zeebe.Group("/order-long-ship")
	orderLongShip.GET("/list", handler.GetOrderLongShipListHandler)
	orderLongShip.GET("/id/:id", handler.GetOrderLongShipHandler)
	orderLongShip.GET("/create-form-data", handler.CreateLongShipFormData)
	orderLongShip.POST("/create", handler.CreateOrderLongShipHandler)
	orderLongShip.GET("/update-form-data/:id", handler.UpdateLongShipFormData)
	orderLongShip.PUT("/update/:id", handler.UpdateOrderLongShipHandler)
	orderLongShip.PUT("/update/load-package", updateOLSLoadPackageHandler)
	orderLongShip.PUT("/update/start-vehicle", updateOLSStartVehicleHandler)
	orderLongShip.PUT("/update/vehicle-arrived", updateOLSVehicleArrivedHandler)
	orderLongShip.PUT("/update/unload-package", updateOLSUnloadPackageHandler)
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
}

func updateOLSLoadPackageHandler(c *gin.Context) {
	userAuthID, err := middleware.GetUserAuthIDInSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler.UpdateOLSLoadPackageHandler(c, userAuthID)
}

func updateOLSStartVehicleHandler(c *gin.Context) {
	userAuthID, err := middleware.GetUserAuthIDInSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler.UpdateOLSStartVehicleHandler(c, userAuthID)
}

func updateOLSVehicleArrivedHandler(c *gin.Context) {
	userAuthID, err := middleware.GetUserAuthIDInSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler.UpdateOLSVehicleArrivedHandler(c, userAuthID)
}

func updateOLSUnloadPackageHandler(c *gin.Context) {
	userAuthID, err := middleware.GetUserAuthIDInSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler.UpdateOLSUnloadPackageHandler(c, userAuthID)
}
