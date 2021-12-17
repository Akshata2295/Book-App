package models

import "time"

type UserRole struct {
	Id   int    `json:"id"`
	Role string `json:"role"`
}

type User struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Email      string `json:"email" gorm:"unique"`
	UserRoleID int    `json:"role_id,string"`
	Password   string `json:"-"`
	Mobile     string `json:""`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}
