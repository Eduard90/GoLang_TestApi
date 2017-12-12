package models

import (
	"time"
	"github.com/google/jsonapi"
	"fmt"
)

type Product struct {
	ID			uint		`gorm:"primary_key" jsonapi:"primary,products"`
	CreatedAt 	time.Time	`jsonapi:"attr,created_at"`
	UpdatedAt 	time.Time	`jsonapi:"attr,updated_at"`
	Name		string		`gorm:"size:100" jsonapi:"attr,name"`
	CategoryID	int			`jsonapi:"attr,category_id"`
	Price		float32		`jsonapi:"attr,price"`
}


func (category Product) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("https://localhost:8080/products/%d", category.ID),
	}
}



