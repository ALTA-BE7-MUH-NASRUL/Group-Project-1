package main

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type User struct {
	gorm.Model
	Name     string `json:"name" form:"name"`
	Email    string `gorm:"unique" json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Phone    string `gorm:"unique" json:"phone" form:"phone"`
	Address  string `json:"address" form:"address"`
	ATM      []ATM
}

type ATM struct {
	gorm.Model
	Number string
	UserID uint
}

func InitDB() {
	connectionString := os.Getenv("group_project")

	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}

func InitialMigration() {
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&ATM{})
}

func init() {
	InitDB()
	InitialMigration()
}

func main() {
	fmt.Println("Hello World!")
}
