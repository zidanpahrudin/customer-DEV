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
       
		&entity.ActivityAttendee{},
		&entity.ActivityCheckin{},
		&entity.ActivityType{},
		&entity.Address{},
		&entity.Contact{},
		&entity.Customer{},
		&entity.Document{},
		&entity.HistoryCustomer{},
		&entity.Event{},
		&entity.EventAttendee{},
		&entity.Group{},
		&entity.Invoice{},
		&entity.Other{},
		&entity.Payment{},
		&entity.Project{},
		&entity.Role{},
		&entity.Sosmed{},
		&entity.Status{},
		&entity.StatusReasons{},
		&entity.Structure{},
		&entity.User{},
		&entity.Stages{},
		&entity.StagesDetail{},
		&entity.Workflows{},
		&entity.WorkflowsDetail{},
		&entity.GroupConfig{},
		&entity.GroupConfigDetail{},
		 &entity.Activity{},
		
		
    }
    for _, t := range tables {
        _ = DB.Migrator().DropTable(t)
    }
}


	// Auto migrate the schema - akan membuat tabel sesuai model Go
	err = DB.AutoMigrate(
		
		&entity.ActivityAttendee{},
		&entity.ActivityCheckin{},
		&entity.ActivityType{},
		&entity.Address{},
		&entity.Contact{},
		&entity.Customer{},
		&entity.Document{},
		&entity.HistoryCustomer{},
		&entity.Event{},
		&entity.EventAttendee{},
		&entity.Group{},
		&entity.Invoice{},
		&entity.Other{},
		&entity.Payment{},
		&entity.Project{},
		&entity.Role{},
		&entity.Sosmed{},
		&entity.Status{},
		&entity.StatusReasons{},
		&entity.Structure{},
		&entity.User{},
		&entity.Stages{},
		&entity.StagesDetail{},
		&entity.Workflows{},
		&entity.WorkflowsDetail{},
		&entity.GroupConfig{},
		&entity.GroupConfigDetail{},
		&entity.Activity{},
		
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
