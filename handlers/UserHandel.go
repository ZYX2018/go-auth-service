package handlers

import (
	"go-auth-service/models"
	"gorm.io/gorm"
	"log"
)

type UserHandel interface {
	CreateUser(user *models.User) (string, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserById(id string) (*models.User, error)
	DeleteUserById(id string) (string, error)
}

type userHandel struct {
	db *gorm.DB
}

func NewUserHandel(db *gorm.DB) UserHandel {
	return &userHandel{db: db}
}

func (handel *userHandel) CreateUser(user *models.User) (string, error) {
	if user == nil {
		return "", nil
	}
	result := handel.db.Create(&user)
	if result.Error != nil {
		log.Fatalln("Error creating user")
		return "", result.Error
	}
	return user.ID, nil
}

func (handel *userHandel) GetUserByUsername(userName string) (*models.User, error) {
	var userModel models.User
	result := handel.db.First(&userModel, "userName = ?", userName)
	if result.Error != nil {
		log.Fatalln("Error getting user")
	}
	return &userModel, result.Error
}

func (handel *userHandel) GetUserById(id string) (*models.User, error) {
	var userModel models.User
	result := handel.db.First(&userModel, "id = ?", id)
	return &userModel, result.Error
}
func (handel *userHandel) DeleteUserById(id string) (string, error) {
	result := handel.db.Delete(&models.User{}, id)
	return id, result.Error
}
