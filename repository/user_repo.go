package repository

import (
	"example/auth/domain"
	"example/auth/mysql"
	"example/auth/usecase"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) usecase.UserRepository {
	return &UserRepo{db: db}
}

func (ur *UserRepo) RegisterUser(user *domain.User) (*domain.User, error) {

	// hash pass menggunakan bcrypt
	fmt.Println(user.Password)
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashPassword)

	fmt.Println(user.Password)
	// insert ke database
	if err := mysql.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("error cause: %d", &err)
	}

	return user, nil
}

func (ur *UserRepo) LoginUser(userName string, password string) (string, error) {

	// ambil data user berdasarkan username
	var user domain.User
	if err := mysql.DB.Where("username = ?", userName).First(&user).Error; err != nil {
		return "", err
	}

	// cek apakah password valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	return user.Username, nil
}

func (ur *UserRepo) ProfilUser() (*[]domain.User, error) {

	var user []domain.User

	if err := mysql.DB.Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil

}
