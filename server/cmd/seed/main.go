package main

import (
	"fmt"
	"log"
	"os"

	"nota-parfume/internal/config"
	"nota-parfume/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	db := config.SetUpDatabaseConnection()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	if err := db.AutoMigrate(&models.Parfume{}, &models.Order{}, &models.OrderItem{}, &models.Admin{}); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}

	// очищаем (только для разработки)
	db.Exec("DELETE FROM order_items")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM parfumes")
	db.Exec("DELETE FROM admins")

	// =====================
	// ADMIN
	// =====================

	hash, err := bcrypt.GenerateFromPassword(
		[]byte("admin12345"),
		bcrypt.DefaultCost,
	)

	if err != nil {
		panic(err)
	}

	admin := models.Admin{
		Name:         "Main Admin",
		Email:        "admin@test.com",
		PasswordHash: string(hash),
		Role:         "admin",
	}

	db.Create(&admin)

	// =====================
	// PARFUMES
	// =====================

	parfumes := []models.Parfume{

		{
			Name:        "Bleu de Chanel",
			Description: "Fresh woody fragrance",
			Brand:       "Chanel",
			Category:    "Men",
			Notes: []string{
				"citrus",
				"cedar",
				"sandalwood",
			},
			PricePerMl: 25,
			AvailableVolumes: []uint{
				30,
				50,
				100,
			},
			ImageUrl: "https://example.com/bleu.jpg",
			IsActive: true,
			Badge:    "Bestseller",
		},

		{
			Name:        "Sauvage",
			Description: "Aromatic fresh scent",
			Brand:       "Dior",
			Category:    "Men",
			Notes: []string{
				"bergamot",
				"pepper",
				"ambroxan",
			},
			PricePerMl: 18,
			AvailableVolumes: []uint{
				30,
				100,
			},
			ImageUrl: "https://example.com/sauvage.jpg",
			IsActive: true,
			Badge:    "Popular",
		},

		{
			Name:        "Black Orchid",
			Description: "Dark floral fragrance",
			Brand:       "Tom Ford",
			Category:    "Unisex",
			Notes: []string{
				"orchid",
				"spices",
				"patchouli",
			},
			PricePerMl: 35,
			AvailableVolumes: []uint{
				50,
				100,
			},
			ImageUrl: "https://example.com/orchid.jpg",
			IsActive: true,
			Badge:    "Premium",
		},
	}

	for _, p := range parfumes {
		db.Create(&p)
	}

	// =====================
	// ORDER
	// =====================

	var bleu models.Parfume

	db.First(
		&bleu,
		"name = ?",
		"Bleu de Chanel",
	)

	order := models.Order{

		CustomerName:    "Ivan",
		ContactMethod:   3,
		Phone:           "+48111111111",
		City:            "Warsaw",
		DeliveryAddress: "Test street 1",
		Status:          "new",
		Comment:         "Call before delivery",

		TotalPrice: 750,

		Items: []models.OrderItem{

			{
				ParfumeID: bleu.ID,
				VolumeMl:  30,
				Quantity:  1,

				PricePerMl: bleu.PricePerMl,
				UnitPrice:  bleu.PricePerMl * 30,
				TotalPrice: bleu.PricePerMl * 30,
			},
		},
	}

	db.Create(&order)

	fmt.Println("Seed completed")

}
