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
	Phone    string `gorm:"unique" json:"phone" form:"phone"`
	Transfer []Transfer
	Top_up   []Top_up
}

type Transfer struct {
	gorm.Model
	UserID         int
	Phone_user     string `gorm:"unique" json:"phone_user" form:"phone_user"`
	Phone_receiver string `gorm:"unique" json:"phone_receiver" form:"phone_receiver"`
	Amount         uint   `json:"amount" form:"amount"`
}

type Top_up struct {
	gorm.Model
	UserID     int
	Phone_user string `gorm:"unique" json:"phone_user" form:"phone_user"`
	Amount     uint   `json:"amount" form:"amount"`
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
	DB.AutoMigrate(&Transfer{})
	DB.AutoMigrate(&Top_up{})
}

func init() {
	InitDB()
	InitialMigration()
}

func main() {

	fmt.Println("Hello World!")
}
