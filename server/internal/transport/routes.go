package transport

import (
	"nota-parfume/internal/middleware"
	"nota-parfume/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	parfuma service.ParfumeService,
	orderService service.OrderService,
	adminService service.AdminService,
	authService service.AuthService,
) {

	authorized := router.Group("")
	authorized.Use(middleware.AuthMiddleware())

	unauthorized := router.Group("")

	parfumeHandler := NewParfumeHandler(parfuma)

	orderHandler := NewOrderHandler(orderService)

	adminHandler := NewAdminHandler(adminService)

	authHandler := NewAuthHandler(authService)

	// auth
	unauthorized.POST("/auth/login", authHandler.Login)

	parfumeHandler.ParfumeRegisterRoutes(authorized, unauthorized)
	orderHandler.OrderRegisterRoutes(authorized, unauthorized)
	adminHandler.AdminRegisterRoutes(authorized, unauthorized)
}
