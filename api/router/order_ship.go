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

	longShip := zeebe.Group("/long-ship")
	longShip.GET("/list", handler.GetLongShipListHandler)
	longShip.GET("/id/:id", handler.GetLongShipHandler)
	longShip.GET("/create-form-data", handler.CreateLongShipFormData)
	longShip.POST("/create", handler.CreateLongShipHandler)
	longShip.GET("/update-form-data/:id", handler.UpdateLongShipFormData)
	longShip.PUT("/update/:id", handler.UpdateLongShipHandler)
	longShip.PUT("/update/load-package", updateLSLoadPackageHandler)
	longShip.PUT("/update/start-vehicle", updateLSStartVehicleHandler)
	longShip.PUT("/update/vehicle-arrived", updateLSVehicleArrivedHandler)
	longShip.PUT("/update/unload-package", updateLSUnloadPackageHandler)
	longShip.DELETE("/delete/:id", handler.DeleteLongShipHandler)

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

func updateLSLoadPackageHandler(c *gin.Context) {
	userAuthID, err := middleware.GetUserAuthIDInSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler.UpdateLSLoadPackageHandler(c, userAuthID)
}

func updateLSStartVehicleHandler(c *gin.Context) {
	userAuthID, err := middleware.GetUserAuthIDInSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler.UpdateLSStartVehicleHandler(c, userAuthID)
}

func updateLSVehicleArrivedHandler(c *gin.Context) {
	userAuthID, err := middleware.GetUserAuthIDInSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler.UpdateLSVehicleArrivedHandler(c, userAuthID)
}

func updateLSUnloadPackageHandler(c *gin.Context) {
	userAuthID, err := middleware.GetUserAuthIDInSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler.UpdateLSUnloadPackageHandler(c, userAuthID)
}
