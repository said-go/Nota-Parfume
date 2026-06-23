package main

import (
	"log"
	"os"

	"nota-parfume/internal/config"
	"nota-parfume/internal/models"
	"nota-parfume/internal/repository"
	"nota-parfume/internal/service"
	"nota-parfume/internal/storage"
	"nota-parfume/internal/transport"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db := config.SetUpDatabaseConnection()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	if err := db.AutoMigrate(&models.Parfume{}, &models.Order{}, &models.OrderItem{}, &models.Admin{}); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}

	parfumeRepo := repository.NewParfumeRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	yandexStorage := storage.NewYandexStorage(os.Getenv("YANDEX_TOKEN"))
	
	parfumeService := service.NewParfumeService(parfumeRepo, yandexStorage)
	orderService := service.NewOrderService(orderRepo, parfumeRepo)
	adminService := service.NewAdminService(adminRepo)
	authService := service.NewAuthService(adminRepo)

	transport.RegisterRoutes(router, parfumeService, orderService, adminService, *authService)

	// fmt.Printf("Starting Nota-Parfume server on %s\n", addr)
	if err := router.Run(); err != nil {
		log.Fatalf("server gin run error: %v", err)
	}
}
