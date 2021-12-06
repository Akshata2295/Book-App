package models

import "time"

type Book struct {
	ID         int    `json:"id"`
	Title      string `json:"title" gorm:"unique"`
	CreatedBy  int
	User       User `gorm:"foreignKey:CreatedBy" json:"-"`
	CategoryId int  `json:"category_id,string"`
	// Category Category `gorm:"foreignKey:CategoryId" json:"-"`
	Price      int  `json:"price"`
	// Meta
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

