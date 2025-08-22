package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"customer-api/internal/entity"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection parameters
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Create connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port)

	// Connect to database
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Assign database connection to global DB variable
	DB = database

	fmt.Println("Starting auto migration...")


	isProd, _ := strconv.ParseBool(os.Getenv("IS_PRODUCTION"))
	// Drop semua tabel yang bermasalah untuk memastikan skema bersih
	// Drop tabel dengan foreign key terlebih dahulu
	if !isProd {
    tables := []interface{}{
        &entity.Status{},
        &entity.Payment{},
        &entity.Invoice{},
        &entity.ActivityCheckin{},
        &entity.ActivityAttendee{},
        &entity.Activity{},
        &entity.User{},
        &entity.Address{},
        &entity.Sosmed{},
        &entity.Contact{},
        &entity.Structure{},
        &entity.Other{},
        &entity.Customer{},
        &entity.Role{},
    }
    for _, t := range tables {
        _ = DB.Migrator().DropTable(t)
    }
}


	// Auto migrate the schema - akan membuat tabel sesuai model Go
	err = DB.AutoMigrate(
		&entity.Role{},
		&entity.User{},
		&entity.Customer{},
		&entity.Address{},
		&entity.Sosmed{},
		&entity.Contact{},
		&entity.Structure{},
		&entity.Group{},
		&entity.Other{},
		&entity.Activity{},
		&entity.ActivityCheckin{},
		&entity.ActivityAttendee{},
		&entity.Invoice{},
		&entity.Payment{},
		&entity.Status{},
		&entity.StatusReasons{},
		&entity.Document{},
		

		
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Insert default roles if they don't exist
	var adminRole entity.Role
	result := DB.Where("role_name = ?", "Admin").First(&adminRole)
	if result.RowsAffected == 0 {
		DB.Create(&entity.Role{RoleName: "Admin"})
		fmt.Println("Created Admin role")
	}

	var userRole entity.Role
	result = DB.Where("role_name = ?", "User").First(&userRole)
	if result.RowsAffected == 0 {
		DB.Create(&entity.Role{RoleName: "User"})
		fmt.Println("Created User role")
	}

	fmt.Println("Database connected and migrated successfully!")
}
