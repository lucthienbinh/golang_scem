package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/handler"
)

// OrderShipRoutes to manage user model
func OrderShipRoutes(rg *gin.RouterGroup) {

	stateService := rg.Group("/state-service")
	if os.Getenv("STATE_SERVICE") == "1" {
		zeebe := stateService.Group("/zeebe")
		zeebe.GET("/deploy-full-ship-workflow", handler.DeployWorkflowFullShipHandlerZB)
		zeebe.GET("/deploy-long-ship-workflow", handler.DeployWorkflowLongShipHandlerZB)
	} else {
		stateScem := stateService.Group("/state-scem")
		stateScem.GET("/deploy-full-ship-workflow", handler.DeployWorkflowFullShipHandlerSS)
		stateScem.GET("/deploy-long-ship-workflow", handler.DeployWorkflowLongShipHandlerSS)
		stateScem.GET("/test-create-instance", handler.CreateWorkflowInstanceSS)
	}

	longShip := rg.Group("/long-ship")
	longShip.GET("/list", handler.GetLongShipListHandler)
	longShip.GET("/id/:id", handler.GetLongShipHandler)
	longShip.GET("/create-form-data", handler.CreateLongShipFormData)
	longShip.POST("/create", handler.CreateLongShipHandler)
	longShip.GET("/update-form-data/:id", handler.UpdateLongShipFormData)
	longShip.PUT("/update/:id", handler.UpdateLongShipHandler)
	longShip.PUT("/update-load-package/:id", updateLSLoadPackageHandler)
	longShip.PUT("/update-start-vehicle/:id", updateLSStartVehicleHandler)
	longShip.PUT("/update-vehicle-arrived/:id", updateLSVehicleArrivedHandler)
	longShip.PUT("/update-unload-package/:id", updateLSUnloadPackageHandler)
	longShip.DELETE("/delete/:id", handler.DeleteLongShipHandler)

	orderShortShip := rg.Group("/order-short-ship")
	orderShortShip.GET("/list", handler.GetOrderShortShipListHandler)
	orderShortShip.GET("/id/:id", handler.GetOrderShortShipHandler)
	orderShortShip.PUT("/update/shipper-called/:id", handler.UpdateOSSShipperCalledHandler)
	orderShortShip.PUT("/update/shipper-received-money/:id", handler.UpdateOSSShipperReceivedMoneyHandler)
	orderShortShip.PUT("/update/shipper-shipped/:id", handler.UpdateOSSShipperShippedHandler)
	orderShortShip.PUT("/update/shipper-confirmed/:id", handler.UpdateOSSShipperConfirmedHandler)
	orderShortShip.PUT("/update/cancel-order/:id", handler.CancelOrderShortShipHandler)
}

// Todo: read user auth id in session instead of hard code
func updateLSLoadPackageHandler(c *gin.Context) {
	// userAuthID, err := middleware.GetUserAuthIDInSession(c)
	// if err != nil {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }
	// handler.UpdateLSLoadPackageHandler(c, userAuthID)
	handler.UpdateLSLoadPackageHandler(c, uint(4))
}

func updateLSStartVehicleHandler(c *gin.Context) {
	// userAuthID, err := middleware.GetUserAuthIDInSession(c)
	// if err != nil {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }
	// handler.UpdateLSStartVehicleHandler(c, userAuthID)
	handler.UpdateLSStartVehicleHandler(c, uint(3))
}

func updateLSVehicleArrivedHandler(c *gin.Context) {
	// userAuthID, err := middleware.GetUserAuthIDInSession(c)
	// if err != nil {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }
	// handler.UpdateLSVehicleArrivedHandler(c, userAuthID)
	handler.UpdateLSVehicleArrivedHandler(c, uint(3))
}

func updateLSUnloadPackageHandler(c *gin.Context) {
	// userAuthID, err := middleware.GetUserAuthIDInSession(c)
	// if err != nil {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }
	// handler.UpdateLSUnloadPackageHandler(c, userAuthID)
	handler.UpdateLSUnloadPackageHandler(c, uint(4))
}
