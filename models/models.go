package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string   `gorm:"unique;not null"`
	Username string   `gorm:"not null"`
	Password string   `gorm:"not null"`
	Roles    string   `gorm:"not null"`
	Profile  *Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID"`
}

type Profile struct {
	gorm.Model
	UserID    uint   `gorm:"unique;not null"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Phone     string `gorm:"not null"`
	Address   string `gorm:"not null"`
	Avatar    string `gorm:"not null"`
}

type Course struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Price       float64 `gorm:"default:0"`
	Thumbnail   string
	UserID      uint
	User        *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID"`
}

type Enrollment struct {
	gorm.Model
	UserID   uint
	User     *User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID"`
	CourseID uint
	Course   *Course `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CourseID"`
}

type Lesson struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
	Video       string
	CourseID    uint
	Course      *Course `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CourseID"`
}

type Quiz struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
	CourseID    uint
	Course      *Course `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CourseID"`
}

type Question struct {
	gorm.Model
	Question string `gorm:"not null"`
	Answer   string `gorm:"not null"`
	QuizID   uint
	Quiz     *Quiz `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:QuizID"`
}

type Answer struct {
	gorm.Model
	Answer     string `gorm:"not null"`
	QuestionID uint
	Question   *Question `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:QuestionID"`
}

type UserCourse struct {
	gorm.Model
	UserID   uint
	User     *User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID"`
	CourseID uint
	Course   *Course `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CourseID"`
}

type UserQuiz struct {
	gorm.Model
	UserID uint
	User   *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID"`
	QuizID uint
	Quiz   *Quiz `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:QuizID"`
}

type UserAnswer struct {
	gorm.Model
	UserID   uint
	User     *User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID"`
	AnswerID uint
	Answer   *Answer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:AnswerID"`
}
