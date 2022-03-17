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
	Name       string `json:"name" form:"name"`
	Phone_user string `gorm:"unique" json:"phone" form:"phone"`
	Balance    uint   `json:"balance" form:"balance"`
	Transfer   []Transfer
	Top_up     []Top_up
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
	fmt.Println("1. Create Account\n 2. Read Account\n 3. Update Account\n 4. Delete Account\n 5. Top-Up\n 6. Transfer\n 7.Read History Top-Up\n 8.Read History Transfer\n Choose menu: ")
	var Choose string
	fmt.Scanln(&Choose)

	switch Choose {
	case "1":
		fmt.Println("Create")
	case "2":
		fmt.Println("Read")
	case "3":
		fmt.Println("Update")
	case "4":
		fmt.Println("Delete")
	case "5":
		fmt.Println("Top-Up")
		top_up := Top_up{}
		user := User{}
		fmt.Println("Insert Your Phone Number: ")
		fmt.Scan(&top_up.Phone_user)
		fmt.Println("Insert Amount: Rp. ")
		fmt.Scan(&top_up.Amount)
		tx := DB.Save(&top_up)
		if tx.Error != nil {
			fmt.Println("error when insert phone number")
		}
		if tx.RowsAffected == 0 {
			fmt.Println("insert failed")
		}

		DB.Model(&user).Select("Balance").Where("Phone_user=?", top_up.Phone_user).Updates(User{Balance: user.Balance + top_up.Amount})

	case "6":
		fmt.Println("Transfer")
		// newTransfer := Transfer{}
		// user := User{}
		// fmt.Println("Insert Your Phone Number: ")
		// fmt.Scan(&newTransfer.Phone_user)
		// fmt.Println("Insert Amount: Rp. ")
		// fmt.Scan(&top_up.Amount)
		// fmt.Println("Insert Destination Phone Number: ")
		// fmt.Scan(&top_up.Phone_receiver)
		// Balance_Amount := DB.Select("balance").Where("phone_user:?", &top_up.Phone_user).Find(&user)
		// if Balance_Amount > &top_up.Amount {
		// }
	case "7":
		fmt.Println("history top-up")
	case "8":
		fmt.Println("history transfer")
	}
}
