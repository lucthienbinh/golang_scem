package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/handler"
)

// WebUserRoutes to manage user model
func WebUserRoutes(rg *gin.RouterGroup) {

	customer := rg.Group("/customer")
	customer.GET("/list", handler.GetCustomerListHandler)
	customer.GET("/id/:id", handler.GetCustomerHandler)
	customer.POST("/create", handler.CreateCustomerHandler)
	customer.PUT("/update/:id", handler.UpdateCustomerHandler)
	customer.DELETE("/delete/:id", handler.DeleteCustomerHandler)

	customerCredit := rg.Group("/customer-credit")
	customerCredit.GET("/list", handler.GetCustomerCreditListHandler)
	customerCredit.PUT("/update/validattion/:id", handler.UpdateCustomerCreditValidationHandler)
	customerCredit.PUT("/update/balance/:id", handler.UpdateCustomerCreditBalanceHandler)

	employee := rg.Group("/employee")
	employee.GET("/list", handler.GetEmployeeListHandler)
	employee.GET("/id/:id", handler.GetEmployeeHandler)
	employee.GET("/create-form-data", handler.CreateEmployeeFormData)
	employee.POST("/create", handler.CreateEmployeeHandler)
	employee.POST("/upload/image", handler.ImageEmployeeHandler)
	employee.GET("/update-form-data/:id", handler.UpdateEmployeeFormData)
	employee.PUT("/update/:id", handler.UpdateEmployeeHandler)
	employee.DELETE("/delete/:id", handler.DeleteEmployeeHandler)

	employeeType := rg.Group("/employee-type")
	employeeType.GET("/list", handler.GetEmployeeTypeListHandler)
	employeeType.GET("/id/:id", handler.GetEmployeeTypeHandler)
	employeeType.POST("/create", handler.CreateEmployeeTypeHandler)
	employeeType.PUT("/update/:id", handler.UpdateEmployeeTypeHandler)
	employeeType.DELETE("/delete/:id", handler.DeleteEmployeeTypeHandler)

	deliveryLocation := rg.Group("/delivery-location")
	deliveryLocation.GET("/list", handler.GetDeliveryLocationListHandler)
	deliveryLocation.GET("/id/:id", handler.GetDeliveryLocationHandler)
	deliveryLocation.POST("/create", handler.CreateDeliveryLocationHandler)
	deliveryLocation.PUT("/update/:id", handler.UpdateDeliveryLocationHandler)
	deliveryLocation.DELETE("/delete/:id", handler.DeleteDeliveryLocationHandler)

}

// AppUserRoutes to manage user model
func AppUserRoutes(rg *gin.RouterGroup) {

	customer := rg.Group("/customer")
	customer.GET("/id/:id", handler.GetCustomerHandler)
	customer.PUT("/update/:id", handler.UpdateCustomerHandler)

	customerNotification := rg.Group("/customer-notification")
	customerNotification.GET("/list/customer-id/:id", handler.GetCustomerNotificationListByCustomerIDHandler)

	employee := rg.Group("/employee")
	employee.GET("/id/:id", handler.GetEmployeeHandler)

}
