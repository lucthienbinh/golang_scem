package router

import (
	"github.com/gin-gonic/gin"
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
	order.POST("/upload/image", handler.ImageOrderHandler)
	order.PUT("/update/:id", handler.UpdateOrderInfoHandler)
	order.DELETE("/delete/:id", handler.DeleteOrderInfoHandler)

	orderPay := rg.Group("/order-pay")
	orderPay.GET("/list", handler.GetOrderPayListHandler)
	orderPay.POST("/create-step-one", handler.CreateOrderPayStepOneHandler)
	orderPay.POST("/create-step-two", handler.CreateOrderPayStepTwoHandler)
	orderPay.PUT("/update-payment-confirm/orderid/:id", handler.UpdateOrderPayConfirmHandler)

	orderVoucher := rg.Group("/order-voucher")
	orderVoucher.GET("/list", handler.GetOrderVoucherListHandler)
}

// AppOrderRoutes to manage order model
func AppOrderRoutes(rg *gin.RouterGroup) {

	order := rg.Group("/order")
	order.GET("/list/customer-id/:id", handler.GetOrderListByCustomerIDHandler)
	order.GET("/id/:id", handler.GetOrderInfoHandler)
	order.POST("/create-use-voucher", handler.CreateOrderUseVoucherInfoHandler)
	order.POST("/upload/image", handler.ImageOrderHandler)

	orderPay := rg.Group("/order-pay")
	orderPay.POST("/create-step-one", handler.CreateOrderPayStepOneHandler)
	orderPay.POST("/create-step-two", handler.CreateOrderPayStepTwoHandler)
	orderPay.PUT("/update-payment-confirm/orderid/:id", handler.UpdateOrderPayConfirmHandler)

	orderVoucher := rg.Group("/order-voucher")
	orderVoucher.GET("/list", handler.GetOrderVoucherListHandler)
}

// Todo: get user token first then go to handler
// func updateOrderPayConfirmHandler(c *gin.Context) {
// 	userAuthID, err := middleware.GetUserAuthIDInToken(c)
// 	if err != nil {
// 		c.AbortWithStatus(http.StatusUnauthorized)
// 		return
// 	}
// 	handler.UpdateOrderPayCustomerConfirmHandler(c, userAuthID)
// }
