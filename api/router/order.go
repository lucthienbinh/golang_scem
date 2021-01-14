package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/api/middleware"
	"github.com/lucthienbinh/golang_scem/internal/handler"
)

// WebOrderRoutes to manage order model
func WebOrderRoutes(rg *gin.RouterGroup) {

	transportType := rg.Group("/transport-type")
	transportType.GET("/list", handler.GetTransportTypeListHandler)
	transportType.GET("/id/:id", handler.GetTransportTypeHandler)
	transportType.POST("/create", handler.CreateTransportTypeHandler)
	transportType.PUT("/update/:id", handler.UpdateTransportTypeHandler)
	transportType.DELETE("/delete/:id", handler.DeleteTransportTypeHandler)

	order := rg.Group("/order")
	order.GET("/list", handler.GetOrderInfoListHandler)
	order.GET("/id/:id", handler.GetOrderInfoHandler)
	order.GET("/create-form-data", handler.CreateOrderFormData)
	order.POST("/create", handler.CreateOrderInfoHandler)
	order.PUT("/update/:id", handler.UpdateOrderInfoHandler)
	order.DELETE("/delete/:id", handler.DeleteOrderInfoHandler)

	orderPay := rg.Group("/order-pay")
	orderPay.GET("/list", handler.GetOrderPayListHandler)
	orderPay.GET("/id/:id", handler.GetOrderPayHandler)
	orderPay.POST("/create-step-one", handler.CreateOrderPayStepOneHandler)
	orderPay.POST("/create-step-two", handler.CreateOrderPayStepTwoHandler)
	orderPay.PUT("/update/customer-confirm/:id", updateOrderPayCustomerConfirmHandler)
	orderPay.PUT("/update/employee-confirm/:id", updateOrderPayEmployeeConfirmHandler)
}

func updateOrderPayCustomerConfirmHandler(c *gin.Context) {
	userAuthID, err := middleware.GetUserAuthIDInToken(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler.UpdateOrderPayCustomerConfirmHandler(c, userAuthID)
}

func updateOrderPayEmployeeConfirmHandler(c *gin.Context) {
	userAuthID, err := middleware.GetUserAuthIDInSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler.UpdateOrderPayEmployeeConfirmHandler(c, userAuthID)
}
