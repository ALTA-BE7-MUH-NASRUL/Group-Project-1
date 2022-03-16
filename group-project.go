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
	Name  string `json:"name" form:"name"`
	Phone string `gorm:"unique" json:"phone" form:"phone"`
}

type Transfer struct {
	gorm.Model
	User_id        []User
	Phone_user     string `gorm:"unique" json:"phone_user" form:"phone_user"`
	Phone_receiver string `gorm:"unique" json:"phone_receiver" form:"phone_receiver"`
	Amount         uint   `json:"Amount" form:"Amount"`
}

type Top_up struct {
	gorm.Model
	User_id    []User
	Phone_user string `gorm:"unique" json:"phone_user" form:"phone_user"`
	Amount     uint   `json:"Amount" form:"Amount"`
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
