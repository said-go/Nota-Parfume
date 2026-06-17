package main

import (
	"fmt"
	"log"

	"nota-parfume/internal/config"
	"nota-parfume/internal/models"
	"nota-parfume/internal/repository"
	"nota-parfume/internal/service"
	"nota-parfume/internal/transport"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db, err := repository.NewDatabase(cfg.DSN())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&models.Admin{}, &models.Parfume{}, &models.Order{}, &models.OrderItem{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	adminRepo := repository.NewAdminRepository(db)
	parfumeRepo := repository.NewParfumeRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	adminService := service.NewAdminService(adminRepo)
	parfumeService := service.NewParfumeService(parfumeRepo)
	orderService := service.NewOrderService(orderRepo, parfumeRepo)

	app := transport.NewApp(parfumeService, orderService, adminService)
	router := gin.Default()
	app.RegisterRoutes(router)

	addr := cfg.AppAddress()
	fmt.Printf("Starting Nota-Parfume server on %s\n", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
