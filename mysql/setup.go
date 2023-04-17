package mysql

import (
	"example/auth/domain"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/go"))
	if err != nil {
		fmt.Println("Gagal koneksi database")
	}

	db.AutoMigrate(&domain.User{})

	DB = db
}
