package repository

import (
	"SMS/model"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//DatabaseRepository interface
type DatabaseRepository interface {
	UserDatabase
}

type database struct {
	connection *gorm.DB
}

var db *gorm.DB

//NewDatabaseRepository construct to create connection with DB with any model
func NewDatabaseRepository() DatabaseRepository {
	godotenv.Load()
	dsn := os.Getenv("DB_LINK")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	} else {
		fmt.Println("connected to database")
		db.Debug().AutoMigrate(&model.User{}, &model.Token{})
	}

	return &database{
		connection: db,
	}
}
