package utils

import (
	"SMS/model"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//DB connection function with database
func DB() *gorm.DB {
	godotenv.Load()
	dsn := os.Getenv("DB_LINK")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	} else {
		fmt.Println("connected to database")
		db.Debug().AutoMigrate(&model.User{}, &model.Token{}, &model.Professor{})
	}
	return db
}
