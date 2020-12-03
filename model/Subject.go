package model

//Subject model
type Subject struct {
	ID       uint   `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Material string `gorm:"type:varchar(255)" json:"material"`
}
