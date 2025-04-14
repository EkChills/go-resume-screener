package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

type AnalyzedResume struct {
	Name	 string `json:"name"`
	Email	 string `json:"email"`
	Phone	 string `json:"phone"`
	Skills	 []string `json:"skills"`
	Education []string `json:"education"`
	Experience []string `json:"experience"`
} 

type Resume struct {
	gorm.Model
	Name string `gorm:"not null" json:"name"`
	Email string `gorm:"not null" json:"email"`
	Phone string `gorm:"not null" json:"phone"`
	Skills string `gorm:"not null" json:"skills"`
	Education string `gorm:"not null" json:"education"`
	Experience string `gorm:"not null" json:"experience"`
	Shortlisted bool `gorm:"default:false" json:"shortlisted"`
	UserID uint `gorm:"not null" json:"user_id"`
}