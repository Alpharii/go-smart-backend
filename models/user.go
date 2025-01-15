package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null"`
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
	Roles    []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	ID          uint        `gorm:"primaryKey"`
	Name        string      `gorm:"unique;not null"`
}
