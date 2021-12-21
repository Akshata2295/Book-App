package models

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-redis/redis"
	//_ "github.com/lib/pq"
)

var DB *gorm.DB
var Rdb *redis.Client

type Settings struct {
	DB_HOST     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
	DB_PORT     string
}

func InitializeSettings() Settings {
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")

	var addr = os.Getenv("REDIS_PORT")
	var pass = os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	switch {
	case DB_HOST == "":
		fmt.Println("Environment variable DB_HOST not set.")
		os.Exit(1)
	case DB_NAME == "":
		fmt.Println("Environment variable DB_NAME not set.")
		os.Exit(1)
	case DB_USER == "":
		fmt.Println("Environment variable DB_USER not set.")
		os.Exit(1)
	case DB_PASSWORD == "":
		fmt.Println("Environment variable DB_PASSWORD not set.")
		os.Exit(1)
	}

	settings := Settings{
		DB_HOST:     DB_HOST,
		DB_NAME:     DB_NAME,
		DB_USER:     DB_USER,
		DB_PASSWORD: DB_PASSWORD,
		DB_PORT:     DB_PORT,
	}

	return settings
}

func createInitialData() {

	var UserRoles = []UserRole{
		{
			Id:   1,
			Role: "admin",
		},
		{
			Id:   2,
			Role: "supervisor",
		},
		{
			Id:   3,
			Role: "customer",
		},
	}

	err := DB.CreateInBatches(UserRoles, 3).Error
	if err != nil {
		fmt.Println("user roles created successfully")
	}
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("SuperPassword@123"), 8)

	err = DB.Create(&User{
		ID:         1,
		FirstName:  "admin",
		LastName:   "admin",
		Email:      "admin@gmail.com",
		UserRoleID: 1,
		Mobile:     "9920319453",
		Password:   string(encryptedPassword),
	}).Error
	if err != nil {
		fmt.Println("Admin already exists")
	}

	err = DB.Create(&User{
		ID:         2,
		FirstName:  "supervisor",
		LastName:   "supervisor",
		Email:      "supervisor@gmail.com",
		UserRoleID: 2,
		Mobile:     "97690101886",
		Password:   string(encryptedPassword),
	}).Error
	if err != nil {
		fmt.Println("Supervisor already exists")
	}

	err = DB.Create(&User{
		ID:         3,
		FirstName:  "customer",
		LastName:   "customer",
		Email:      "customer@gmail.com",
		UserRoleID: 3,
		Mobile:     "9820332549",
		Password:   string(encryptedPassword),
	}).Error
	if err != nil {
		fmt.Println("Customer already exists")
	}
}

func ConnectDataBase() {
	// database := fmt.Sprintf("host=localhost port=5432 user=postgres dbname=postgres password=admin sslmode=disable")
	//database := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_NAME, DB_PASSWORD)
	settings := InitializeSettings()
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", settings.DB_USER, settings.DB_PASSWORD, settings.DB_HOST, settings.DB_PORT, settings.DB_NAME)
	//database := "host=" + settings.DB_HOST + " user=" + settings.DB_USER + " password=" + settings.DB_PASSWORD + " dbname=" + settings.DB_NAME + " port=" + settings.DB_PORT + " sslmode=disable"
	//fmt.Println("conname is\t", database)
	connection, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}
	fmt.Println("Migrating tables")

	

	//normal migration
	err = connection.AutoMigrate(&UserRole{}, &User{}, &Book{}, &Cart{}, &Category{}, &RequestLog{}, &FailedRequestLog{})
	if err != nil {
		fmt.Println("error in migration: ", err)
	}
	fmt.Println("Done migrating")

	//DB.Migrator().CreateTable(&Category{})
	DB = connection

	fmt.Println("Done migrating")
	//connection.AutoMigrate(&User{})

	createInitialData()
}