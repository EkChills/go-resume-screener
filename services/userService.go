package services

import (
	"fmt"

	"github.com/ekchills/go-resume-screener/auth"
	"github.com/ekchills/go-resume-screener/models"
	"gorm.io/gorm"
)

type UserService struct {
	Db *gorm.DB
}

func (u *UserService) Register(email string, password string) error {
	hp := auth.Password{Password: password}
	hashedPassword, err := hp.HashPassword()
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return err
	}
	err = u.Db.Create(&models.User{Email: email, Password: hashedPassword}).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) Login(email string, password string) (*models.User, error) {
	user := models.User{}
	err := u.Db.Where("email = ?", email).First(&user).Error
	if err != nil {
		fmt.Println("Error fetching user:", err)
		return nil, err
	}
	p := auth.Password{Password: password}
	_, err = p.ComparePasswordHash(user.Password)
	if err != nil {
		fmt.Println("Error comparing password:", err, user.Password, password)
		return nil, err
	}
	
	return &user, nil
}
