package models

import (
	"time"
	"github.com/google/jsonapi"
	"fmt"
)

type Category struct {
	ID        	uint 		`gorm:"primary_key" jsonapi:"primary,categories"`
	CreatedAt 	time.Time	`jsonapi:"attr,created_at"`
	UpdatedAt 	time.Time	`jsonapi:"attr,updated_at"`
	Name		string		`gorm:"size:100" jsonapi:"attr,name"`
	Lft			int			`jsonapi:"attr,lft"`
	Rgt			int			`jsonapi:"attr,rgt"`
	Level		int			`jsonapi:"attr,level"`
	ParentID	int			`jsonapi:"attr,parent_id"`
}


func (category Category) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("https://localhost:8080/categories/%d", category.ID),
	}
}