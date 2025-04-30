package user

import (
	"fmt"

	"task_management/internal/util"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *UserRepo
}

func NewUserService(repo *UserRepo) *UserService {
	return &UserService{userRepo: repo}
}

// add user
func (s *UserService) AddUser(userData *AddUserBody, dbConn *gorm.DB) error {

	// check if the user email exists
	if user, err := s.userRepo.GetUserByEmail(userData.Email, dbConn); err != nil || user != nil {
		if err != nil {
			return fmt.Errorf("error fetching user: %w", err)
		}
		return fmt.Errorf("user with email %s already exists", userData.Email)
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// create the user
	userData.Password = string(hashedPassword)

	// Insert the user in the database
	err = s.userRepo.CreateUser(userData, dbConn)
	if err != nil {
		return err
	}

	return nil

}

func (s *UserService) LoginUser(userData LoginUserBody, dbConn *gorm.DB) (string, error) {

	// check if the user email exists
	userInfo, err := s.userRepo.GetUserByEmail(userData.Email, dbConn)
	if err != nil || userInfo == nil {
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("user with email %s does not exist", userData.Email)
	}

	// check user password
	if !util.CheckPasswordHash(userData.Password, userInfo.Password) {
		return "", fmt.Errorf("invalid password")
	}

	// create the access token
	accessToken, err := util.CreateToken(userInfo.Id)
	if err != nil {
		return "", err
	}

	return accessToken, nil

}
