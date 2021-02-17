package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/handler"
	ServiceCommon "github.com/lucthienbinh/golang_scem/internal/service/common"
)

// WebOrderShipRoutes to manage user model
func WebOrderShipRoutes(rg *gin.RouterGroup) {

	stateService := rg.Group("/state-service")
	if os.Getenv("STATE_SERVICE") == "1" {
		zeebe := stateService.Group("/zeebe")
		zeebe.GET("/full-ship-workflow/deploy", ServiceCommon.DeployWorkflowFullShipHandlerZB)
		zeebe.GET("/full-ship-workflow/create-instance", ServiceCommon.CreateInstanceWorkflowFSHandlerZB)
		zeebe.GET("/full-ship-workflow/create-instance-internal-bug", ServiceCommon.CreateInstanceInternalBugWorkflowFSHandlerZB)
		zeebe.GET("/full-ship-workflow/create-instance-missing-param-bug", ServiceCommon.CreateInstanceMissingParamBugWorkflowFSHandlerZB)
		zeebe.GET("/long-ship-workflow/deploy", ServiceCommon.DeployWorkflowLongShipHandlerZB)
	} else {
		stateScem := stateService.Group("/state-scem")
		stateScem.GET("/full-ship-workflow/deploy", ServiceCommon.DeployWorkflowFullShipHandlerSS)
		stateScem.GET("/long-ship-workflow/deploy", ServiceCommon.DeployWorkflowLongShipHandlerSS)
	}

	longShip := rg.Group("/long-ship")
	longShip.GET("/list", handler.GetLongShipListHandler)
	longShip.GET("/id/:id", handler.GetLongShipHandler)
	longShip.GET("/create-form-data", handler.CreateLongShipFormData)
	longShip.POST("/create", handler.CreateLongShipHandler)
	longShip.GET("/update-form-data/:id", handler.UpdateLongShipFormData)
	longShip.PUT("/update/:id", handler.UpdateLongShipHandler)
	longShip.PUT("/update-load-package/qrcode/:qrcode", updateLSLoadPackageHandler)
	longShip.PUT("/update-start-vehicle/qrcode/:qrcode", updateLSStartVehicleHandler)
	longShip.PUT("/update-vehicle-arrived/qrcode/:qrcode", updateLSVehicleArrivedHandler)
	longShip.PUT("/update-unload-package/qrcode/:qrcode", updateLSUnloadPackageHandler)
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

// AppOrderShipRoutes to manage user model
func AppOrderShipRoutes(rg *gin.RouterGroup) {

	longShip := rg.Group("/long-ship")
	longShip.GET("/id/:id", handler.GetLongShipHandler)
	longShip.PUT("/update-load-package/qrcode/:qrcode", updateLSLoadPackageHandler)
	longShip.PUT("/update-start-vehicle/qrcode/:qrcode", updateLSStartVehicleHandler)
	longShip.PUT("/update-vehicle-arrived/qrcode/:qrcode", updateLSVehicleArrivedHandler)
	longShip.PUT("/update-unload-package/qrcode/:qrcode", updateLSUnloadPackageHandler)

	orderShortShip := rg.Group("/order-short-ship")
	orderShortShip.GET("/list/employee-id/:id", handler.GetOrderShortShipListByEmployeeIDHandler)
	orderShortShip.GET("/id/:id", handler.GetOrderShortShipHandler)
	orderShortShip.PUT("/update/shipper-called/:id", handler.UpdateOSSShipperCalledHandler)
	orderShortShip.PUT("/update/shipper-received-money/:id", handler.UpdateOSSShipperReceivedMoneyHandler)
	orderShortShip.PUT("/update/shipper-shipped/:id", handler.UpdateOSSShipperShippedHandler)
	orderShortShip.PUT("/update/shipper-confirmed/:id", handler.UpdateOSSShipperConfirmedHandler)
	orderShortShip.PUT("/update/cancel-order/:id", handler.CancelOrderShortShipHandler)

}
