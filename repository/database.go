package repository

import (
	"SMS/utils"

	"gorm.io/gorm"
)

//DatabaseRepository interface
type DatabaseRepository interface {
	UserDatabase
	ProfessorDatabase
}

type database struct {
	connection *gorm.DB
}

var db *gorm.DB

//NewDatabaseRepository construct to create connection with DB with any model
func NewDatabaseRepository() DatabaseRepository {
	db := utils.DB()
	return &database{
		connection: db,
	}
}
