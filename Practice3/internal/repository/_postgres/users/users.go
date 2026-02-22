package users

import (
	"errors"
	"fmt"
	"time"

	_postgres "Practice3/internal/repository/_postgres"
	"Practice3/pkg/modules"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: time.Second * 5,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT id, name, email, age, is_employed FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user, "SELECT id, name, email, age, is_employed FROM users WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return &user, nil
}
func (r *Repository) CreateUser(user modules.User) (int, error) {
	var id int
	query := `INSERT INTO users (name, email, age, is_employed) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.DB.QueryRow(query, user.Name, user.Email, user.Age, user.IsEmployed).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Failed to create user: %v", err)

	}
	return id, nil
}

func (r *Repository) UpdateUser(id int, user modules.User) error {
	query := `UPDATE users SET name=$1, email=$2, age=$3, is_employed=$4 WHERE id=$5`
	result, err := r.db.DB.Exec(query, user.Name, user.Email, user.Age, user.IsEmployed, id)
	if err != nil {
		return fmt.Errorf("Failed to update user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("users not found: no rows were updated")
	}
	return nil
}

func (r *Repository) DeleteUser(id int) (int64, error) {
	result, err := r.db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return 0, fmt.Errorf("Failed to delete user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, fmt.Errorf("user with id %d not found", id)
	}
	return rowsAffected, nil
}
