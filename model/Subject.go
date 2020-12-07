package model

import (
	"html"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

//Subject model
type Subject struct {
	ID       uint   `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	UserID   uint   `json:"user_id"`
	User     *User  `json:"user,omitempty" gorm:"foreignKey:UserID; constraint:OnDelete:SET NULL;"`
	Material string `gorm:"type:varchar(255)" json:"material"`
}

//Prepare Function to escape strings and auto creating dates
func (subject *Subject) Prepare() {
	subject.Name = html.EscapeString(strings.TrimSpace(subject.Name))
}

//AfterFind hook to change the file url in the response
func (subject *Subject) AfterFind(tx *gorm.DB) (err error) {
	godotenv.Load()
	subject.Material = os.Getenv("AWS_URL") + "/" + subject.Material
	return nil
}
