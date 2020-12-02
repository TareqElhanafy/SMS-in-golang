package model

//Professor model
type Professor struct {
	ID     uint `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	UserID uint
	Age    string `gorm:"type:integer" json:"age"`
	Phone  string `gorm:"type:varchar(20)" json:"phone"`
}
