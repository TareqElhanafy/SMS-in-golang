package model

import (
	"html"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//User Model
type User struct {
	ID        uint      `gorm:"primaryKey;  autoIncrement; not null" json:"id"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	Email     string    `gorm:"type:varchar(100); unique" binding:"required" json:"email"`
	Password  string    `gorm:"varchar(255)" json:"-"`
	Tokens    []string  `json:"tokens" gorm:"type:string"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP" `
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

//BeforeSave Hook function to hash passwords before saveing
func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err1 := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err1 != nil {
		return err1
	}
	u.Password = string(hashedPassword)
	u.UpdatedAt = time.Now()
	return nil
}

//Prepare Function to escape strings and auto creating dates
func (u *User) Prepare() {
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

//GenerateToken function to create token
func (u *User) GenerateToken(ID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	godotenv.Load()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
