package repo

import (
	"fmt"

	"Practice7/internal/entity"
	"Practice7/pkg/postgres"

	"github.com/google/uuid"
)

type UserRepo struct {
	PG *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (u *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	user.ID = uuid.New()
	if err := u.PG.Conn.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) LoginUser(user *entity.LoginUserDTO) (*entity.User, error) {
	var userFromDB entity.User
	if err := u.PG.Conn.Where("username = ?", user.Username).First(&userFromDB).Error; err != nil {
		return nil, fmt.Errorf("Username Not Found: %v", err)
	}
	return &userFromDB, nil
}

func (u *UserRepo) GetUserByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := u.PG.Conn.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return &user, nil
}

func (u *UserRepo) PromoteUser(id uuid.UUID) error {
	result := u.PG.Conn.Model(&entity.User{}).Where("id = ?", id).Update("role", "admin")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}
	return nil
}
