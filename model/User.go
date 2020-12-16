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
	ID        uint       `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	Name      string     `gorm:"type:varchar(100)" json:"name"`
	Email     string     `gorm:"type:varchar(100); unique" json:"email"`
	Password  string     `gorm:"varchar(255)" json:"-"`
	Role      string     `gorm:"type:ENUM('superAdmin','professor')" json:"role"`
	Tokens    []Token    `json:"-" gorm:"foreignKey:UserID; constraint:OnDelete;"`
	Professor *Professor `json:"professor,omitempty" gorm:"foreignKey:UserID; constraint:OnDelete;"`
	Subjects  []Subject  `json:"subjects,omitempty" gorm:"foreignKey:UserID; constraint:OnDelete;"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP" `
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

//Token Model
type Token struct {
	ID     uint   `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	UserID uint   `json:"-" gorm:"not null"`
	Token  string `json:"token" gorm:"type:varchar(255)"`
}

//BeforeSave Hook function to hash passwords before saving
func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.UpdatedAt = time.Now()
	return nil
}

//Prepare Function to escape strings and auto creating dates
func (u *User) Prepare() {
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Role = html.EscapeString(strings.TrimSpace(u.Role))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

//GenerateToken function to create token
func (u *User) GenerateToken(ID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	godotenv.Load()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
