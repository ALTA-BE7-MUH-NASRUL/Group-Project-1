package main

import (
	"errors"
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
	connection := os.Getenv("group_project")

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
	fmt.Println("Enter your choice?")
	fmt.Println("1: Create account")
	fmt.Println("2: Showing accounts")
	fmt.Println("3: Update your username")
	fmt.Println("4: Delete your account)")
	fmt.Println("5: Top-up balance")
	fmt.Println("6: Transfer balance")
	fmt.Println("7: Showing history top-up")
	fmt.Println("8: Showing transfer")

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
		fmt.Println("successfully created, your account is now exist")

	case "2":
		var users []User
		tx := DB.Find(&users)
		if tx.Error != nil {
			fmt.Println("error ", tx.Error)
		}
		for _, value := range users {
			fmt.Println("ID: ", value.ID)
			fmt.Println("Name: ", value.Name)
			fmt.Println("Phone number: ", value.Phone_user)
			fmt.Println("Balance", value.Balance)
			fmt.Println(" ")
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
		fmt.Println("Successfully update your name")

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
		fmt.Println("Top up Balance")
		top_up := Top_up{}
		user := User{}
		fmt.Print("Insert your phone number: ")
		fmt.Scanln(&user.Phone_user)
		fmt.Print("Insert amount: Rp. ")
		fmt.Scanln(&top_up.Amount)
		DB.Where("Phone_user=?", user.Phone_user).First(&user)
		user.Balance = user.Balance + top_up.Amount
		DB.Save(&user)
		top_up.UserID = user.ID
		DB.Create(&top_up)
		fmt.Println("Successfully top up your balance")
	case "6":
		fmt.Println("Transfer Balance")
		transfer := Transfer{}
		user := User{}
		receiver := User{}
		fmt.Print("Insert your phone number: ")
		fmt.Scanln(&user.Phone_user)
		fmt.Print("Insert receiver phone number: ")
		fmt.Scanln(&receiver.Phone_user)
		fmt.Print("Insert amount: Rp. ")
		fmt.Scanln(&transfer.Amount)
		DB.Where("Phone_user=?", user.Phone_user).First(&user)
		DB.Where("Phone_user=?", receiver.Phone_user).First(&receiver)
		if user.Balance < transfer.Amount {
			fmt.Println("Your balance is not enough")
		} else {
			user.Balance = user.Balance - transfer.Amount
			DB.Save(&user)
			receiver.Balance = receiver.Balance + transfer.Amount
			DB.Save(&receiver)
			transfer.UserID = user.ID
			transfer.ReceiverID = receiver.ID
			DB.Create(&transfer)
			fmt.Println("Successfully transfered")
		}
	case "7":
		user := User{}
		fmt.Print("Insert Your Phone Number: ")
		fmt.Scanln(&user.Phone_user)
		fmt.Println("History Top Up")
		err := DB.Preload("Top_up").Where("phone_user = ?", user.Phone_user).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Wrong Phone Number")
		}
		fmt.Println("Name :", user.Name)
		fmt.Println("Current Balance :", user.Balance)
		if len(user.Top_up) == 0 {
			fmt.Println("History do not found")
		}
		for i := len(user.Top_up) - 1; i >= 0; i-- {
			fmt.Println("Amount :", user.Top_up[i].Amount, "\t", "Date :", user.Top_up[i].CreatedAt)
		}
	case "8":
		user := User{}
		fmt.Print("Insert Your Phone Number: ")
		fmt.Scanln(&user.Phone_user)
		fmt.Println("History Transfer")
		err := DB.Preload("Receiver").Where("phone_user = ?", user.Phone_user).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Wrong Phone Number")
		}
		fmt.Println("Name :", user.Name)
		fmt.Println("Current Balance :", user.Balance)
		if len(user.Receiver) == 0 {
			fmt.Println("History do not found")
		}
		for i := len(user.Receiver) - 1; i >= 0; i-- {
			receiver := User{}
			DB.Find(&receiver, user.Receiver[i].ReceiverID)
			fmt.Println("Receiver :", receiver.Name)
			fmt.Println("Amount :", user.Receiver[i].Amount)
			fmt.Println("Date :", user.Receiver[i].CreatedAt)
		}

	}
}
