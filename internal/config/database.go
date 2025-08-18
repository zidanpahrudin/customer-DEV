package config

import (
	"fmt"
	"log"
	"os"

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

	// HAPUS SEMUA BARIS DropTable BERIKUT INI:
	// DB.Migrator().DropTable(&entity.ActivityCheckin{})
	// DB.Migrator().DropTable(&entity.ActivityAttendee{})
	// DB.Migrator().DropTable(&entity.Activity{})
	// DB.Migrator().DropTable(&entity.User{})
	// ... dst (semua sudah dikomentari dengan benar)
	// DB.Migrator().DropTable(&entity.Address{})
	// DB.Migrator().DropTable(&entity.Sosmed{})
	// DB.Migrator().DropTable(&entity.Contact{})
	// DB.Migrator().DropTable(&entity.Structure{})
	// DB.Migrator().DropTable(&entity.Other{})
	// DB.Migrator().DropTable(&entity.Customer{})
	// DB.Migrator().DropTable(&entity.Role{})

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
