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
	Name       string     `json:"name" form:"name"`
	Phone_user string     `gorm:"unique" json:"phone" form:"phone"`
	Balance    uint       `json:"balance" form:"balance"`
	Transfer   []Transfer `gorm:"foreignKey:ReceiverID;reference:ID"`
	Receiver   []Transfer `gorm:"foreignKey:UserID;reference:ID"`
	Top_up     []Top_up   `gorm:"foreignKey:UserID;reference:ID"`
}

type Transfer struct {
	gorm.Model
	UserID     uint
	ReceiverID uint
	Amount     uint `json:"amount" form:"amount"`
}

type Top_up struct {
	gorm.Model
	UserID uint
	Amount uint `json:"amount" form:"amount"`
}

func InitDB() {
	connection := "root:$10Milyar@tcp(localhost:3306)/group1?charset=utf8&parseTime=True&loc=Local"

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
		fmt.Scanln(&newUser.Phone_user)

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
			fmt.Println(value.ID, "-", value.Phone_user, "-", value.Name, "-", value.Balance)
		}
	case "3":
		fmt.Println("Change your name")
		var id uint
		var name string
		fmt.Println("Input your id: ")
		fmt.Scanln(&id)
		fmt.Println("Input your new name: ")
		fmt.Scanln(&name)

		tx := DB.Model(&User{}).Where("id = ?", id).Update("Name", name)
		if tx.Error != nil {

			fmt.Println("error when update data")
		}

	case "4":
		fmt.Println("Delete your account")
		var id int
		fmt.Println("Input your id: ")
		fmt.Scanln(&id)
		tx := DB.Delete(&User{}, id)
		if tx.Error != nil {
			fmt.Println("error when delete data")
		}
		fmt.Println("successfully deleted")
	case "5":
		top_up := Top_up{}
		user := User{}
		fmt.Print("Insert Your Phone Number: ")
		fmt.Scanln(&user.Phone_user)
		fmt.Print("Insert Amount: Rp. ")
		fmt.Scanln(&top_up.Amount)
		DB.Where("Phone_user=?", user.Phone_user).First(&user)
		user.Balance = user.Balance + top_up.Amount
		DB.Save(&user)
		top_up.UserID = user.ID
		DB.Create(&top_up)
		fmt.Println("Transaksi Berhasil")
	case "6":
		transfer := Transfer{}
		user := User{}
		receiver := User{}
		fmt.Print("Insert Your Phone Number: ")
		fmt.Scanln(&user.Phone_user)
		fmt.Print("Insert Destination Phone Number: ")
		fmt.Scanln(&receiver.Phone_user)
		fmt.Print("Insert Amount: Rp. ")
		fmt.Scanln(&transfer.Amount)
		DB.Where("Phone_user=?", user.Phone_user).First(&user)
		DB.Where("Phone_user=?", receiver.Phone_user).First(&receiver)
		if user.Balance < transfer.Amount {
			fmt.Println("Your Balance is Not Enough")
		} else {
			user.Balance = user.Balance - transfer.Amount
			DB.Save(&user)
			receiver.Balance = receiver.Balance + transfer.Amount
			DB.Save(&receiver)
			transfer.UserID = user.ID
			transfer.ReceiverID = receiver.ID
			DB.Create(&transfer)
			fmt.Println("Transaksi Berhasil")
		}
	case "7":
		fmt.Println("history top-up")
		var top_up []Top_up
		tx := DB.Find(&top_up)
		if tx.Error != nil {
			fmt.Println("error ", tx.Error)
		}
		for _, value := range top_up {
			fmt.Println(value.ID, "-", value.Amount, "-", value.UpdatedAt, "-", value.CreatedAt)
		}

	case "8":
		fmt.Println("history transfer")
		var transfer []Transfer
		tx := DB.Find(&transfer)
		if tx.Error != nil {
			fmt.Println("error ", tx.Error)
		}
		for _, value := range transfer {
			fmt.Println(value.ID, "-", value.Amount, "-", value.UpdatedAt, "-", value.CreatedAt)
		}
	}
}
