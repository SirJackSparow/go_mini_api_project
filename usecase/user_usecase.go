package usecase

import (
	"example/auth/domain"
)

type UserRepository interface {
	RegisterUser(user *domain.User) (*domain.User, error)
	LoginUser(username string, password string) (string, error)
	ProfilUser() (*[]domain.User, error)
}

type UserUseCase struct {
	userRepo UserRepository
}

func NewUserUseCase(userRepo UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uuc *UserUseCase) RegisterUser(user *domain.User) (*domain.User, error) {
	return uuc.userRepo.RegisterUser(user)
}

func (uuc *UserUseCase) LoginUser(userName string, password string) (string, error) {
	return uuc.userRepo.LoginUser(userName, password)
}

func (uuc *UserUseCase) ProfilUser() (*[]domain.User, error) {
	return uuc.userRepo.ProfilUser()
}
