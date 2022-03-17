package main

import (
	"fmt"

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
	connection := "root:$10Milyar@tcp(localhost:3306)/GroupProject?charset=utf8&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(connection), &gorm.Config{})

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
	fmt.Println("Masukkan pilihan anda? (1: create account)/(2: read your account)/(3: update your account)/(4: delete your account)/(5: top-up balance)/(6: transfer balance)/(7: history top-up)/(8: history transfer)")
	var pilihan string
	fmt.Scanln(&pilihan)

	switch pilihan {
	case "1":

		newUser := User{}
		fmt.Println("Enter your name:")
		fmt.Scanln(&newUser.Name)
		fmt.Println("Enter your phone number:")
		fmt.Scanln(&newUser.Phone)

		tx := DB.Save(&newUser)
		if tx.Error != nil {

			fmt.Println("error when insert data")
		}
		if tx.RowsAffected == 0 {
			fmt.Println("insert failed")
		}
		fmt.Println("successfully created")

	case "2":
		var users []User
		tx := DB.Find(&users)
		if tx.Error != nil {
			fmt.Println("error ", tx.Error)
		}
		for _, value := range users {
			fmt.Println(value.ID, "-", value.Name)
		}
	case "3":
		var id int
		var name string
		fmt.Println("Input your id: ")
		fmt.Scanln(&id)
		fmt.Println("Input your new name: ")
		fmt.Scanln(&name)

		tx := DB.Model(&User{}).Where("id = ?", id).Update("Name", name)
		if tx.Error != nil {

			fmt.Println("error when update data")
		}

	}
}
