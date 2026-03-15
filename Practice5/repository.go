package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPaginatedUsers(filter UserFilter) (PaginatedResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 5
	}

	offset := (filter.Page - 1) * filter.PageSize

	whereParts := []string{"1=1"}
	args := make([]any, 0)
	argPos := 1

	if filter.ID != nil {
		whereParts = append(whereParts, fmt.Sprintf("id = $%d", argPos))
		args = append(args, *filter.ID)
		argPos++
	}

	if filter.Name != nil {
		whereParts = append(whereParts, fmt.Sprintf("name ILIKE $%d", argPos))
		args = append(args, "%"+*filter.Name+"%")
		argPos++
	}

	if filter.Email != nil {
		whereParts = append(whereParts, fmt.Sprintf("email ILIKE $%d", argPos))
		args = append(args, "%"+*filter.Email+"%")
		argPos++
	}

	if filter.Gender != nil {
		whereParts = append(whereParts, fmt.Sprintf("gender ILIKE $%d", argPos))
		args = append(args, *filter.Gender)
		argPos++
	}

	if filter.BirthDate != nil {
		whereParts = append(whereParts, fmt.Sprintf("birth_date = $%d::date", argPos))
		args = append(args, *filter.BirthDate)
		argPos++
	}

	whereClause := strings.Join(whereParts, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", whereClause)

	var totalCount int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return PaginatedResponse{}, fmt.Errorf("count users: %w", err)
	}

	orderBy := sanitizeOrderBy(filter.OrderBy)

	query := fmt.Sprintf(`
		SELECT id, name, email, gender, birth_date
		FROM users
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argPos, argPos+1)

	dataArgs := append(args, filter.PageSize, offset)

	rows, err := r.db.Query(query, dataArgs...)
	if err != nil {
		return PaginatedResponse{}, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Gender, &user.BirthDate); err != nil {
			return PaginatedResponse{}, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return PaginatedResponse{}, fmt.Errorf("rows users: %w", err)
	}

	return PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}, nil
}

func (r *Repository) GetCommonFriends(user1ID, user2ID int) ([]User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM user_friends uf1
		JOIN user_friends uf2 ON uf1.friend_id = uf2.friend_id
		JOIN users u ON u.id = uf1.friend_id
		WHERE uf1.user_id = $1 AND uf2.user_id = $2
		ORDER BY u.id
	`

	rows, err := r.db.Query(query, user1ID, user2ID)
	if err != nil {
		return nil, fmt.Errorf("query common friends: %w", err)
	}
	defer rows.Close()

	friends := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Gender, &user.BirthDate); err != nil {
			return nil, fmt.Errorf("scan common friend: %w", err)
		}
		friends = append(friends, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows common friends: %w", err)
	}

	return friends, nil
}

func sanitizeOrderBy(orderBy string) string {
	orderBy = strings.TrimSpace(orderBy)
	if orderBy == "" {
		return "id ASC"
	}

	direction := "ASC"
	if strings.HasPrefix(orderBy, "-") {
		direction = "DESC"
		orderBy = strings.TrimPrefix(orderBy, "-")
	}

	allowedColumns := map[string]string{
		"id":         "id",
		"name":       "name",
		"email":      "email",
		"gender":     "gender",
		"birth_date": "birth_date",
	}

	column, ok := allowedColumns[orderBy]
	if !ok {
		return "id ASC"
	}

	return column + " " + direction
}
