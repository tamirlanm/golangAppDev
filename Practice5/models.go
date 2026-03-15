package main

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}

type UserFilter struct {
	ID        *int
	Name      *string
	Email     *string
	Gender    *string
	BirthDate *time.Time
	OrderBy   string
	Page      int
	PageSize  int
}

type ErrorResponse struct {
	Error string `json:"error"`
}
