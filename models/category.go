package models

import "time"

type Category struct {
	ID           int    `json:"id"`
	CategoryName string `json:"category_name" gorm:"unique"`
	CreatedBy    int
	User         User `gorm:"foreignKey:CreatedBy" json:"-"`
	// Meta
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
