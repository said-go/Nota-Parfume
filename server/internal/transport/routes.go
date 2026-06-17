package transport

import (
	"nota-parfume/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	parfuma service.ParfumeService,
	orderService service.OrderService,
	adminService service.AdminService,
) {

	authorized := router.Group("")
	// authorized.Use(middleware.InternalAuthMiddleware())

	unauthorized := router.Group("")
	
	parfumeHandler := NewParfumeHandler(parfuma)
	
	orderHandler := NewOrderHandler(orderService)
	
	adminHandler := NewAdminHandler(adminService)
	
	parfumeHandler.ParfumeRegisterRoutes(authorized, unauthorized)
	orderHandler.OrderRegisterRoutes(authorized, unauthorized)
	adminHandler.AdminRegisterRoutes(authorized, unauthorized)
}
