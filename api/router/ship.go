package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/handler"
)

// ShipRoutes to manage ship model
func ShipRoutes(rg *gin.RouterGroup) {

	orderLongShip := rg.Group("/order-long-ship")
	orderLongShip.GET("/list", handler.GetOrderLongShipListHandler)
	orderLongShip.GET("/id/:id", handler.GetOrderLongShipHandler)
	orderLongShip.PUT("/update/load-package", handler.UpdateOLSLoadPackageHandler)
	orderLongShip.PUT("/update/start-vehicle", handler.UpdateOLSStartVehicleHandler)
	orderLongShip.PUT("/update/vehicle-arrived", handler.UpdateOLSVehicleArrivedHandler)
	orderLongShip.PUT("/update/unload-package", handler.UpdateOLSUnloadPackageHandler)
	orderLongShip.DELETE("/delete/:id", handler.DeleteOrderLongShipHandler)

	orderShortShip := rg.Group("/order-short-ship")
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
