package user

import (
	"fmt"
	"log"
	"task_management/pkg/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

type UserLoginResponse struct {
	Id       int
	Password string
}

func (u *UserRepo) GetUserByEmail(email string, db *gorm.DB) (*UserLoginResponse, error) {
	var user models.User

	result := db.Select("*").
		Where("email = ?", email).
		Table("users").
		Scan(&user)

	if result.RowsAffected == 0 {
		error := fmt.Errorf("no user found with email %v", email)
		log.Println(error)
		return nil, nil
	}

	// some other error occurred
	if result.Error != nil {
		error := fmt.Errorf("error while getting user details for %v: %v", email, result.Error)
		log.Println(error)
		return nil, error
	}

	return &UserLoginResponse{
		Id:       int(user.ID),
		Password: user.Password,
	}, nil
}

func (u *UserRepo) CreateUser(user *AddUserBody, db *gorm.DB) error {

	// convert the user to models.User
	userModel := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	// Insert the user in the database
	result := db.Create(&userModel)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	return nil
}
