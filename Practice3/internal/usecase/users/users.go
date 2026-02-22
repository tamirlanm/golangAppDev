package usecase

import (
	"Practice3/internal/repository"
	"Practice3/pkg/modules"
	"fmt"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUsers() ([]modules.User, error) {
	return u.repo.GetUsers()
}

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUsecase) CreateUser(user modules.User) (int, error) {
	return u.repo.CreateUser(user)
}

func (u *UserUsecase) UpdateUser(id int, user modules.User) error {
	return u.repo.UpdateUser(id, user)
}

func (u *UserUsecase) DeleteUser(id int) (int64, error) {
	return u.repo.DeleteUser(id)
}

func (u *UserUsecase) Healthcheck() string {
	return fmt.Sprintf("OK")
}
