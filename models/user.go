package models

import (
	"time"
)


type User struct {
	ID			uint		`gorm:"primary_key" jsonapi:"primary,users"`
	CreatedAt 	time.Time	`jsonapi:"attr,created_at"`
	UpdatedAt 	time.Time	`jsonapi:"attr,updated_at"`
	Login		string		`gorm:"size:100" jsonapi:"attr,login"`
	Password	string		`gorm:"size:255" jsonapi:"attr,login"`
	IsAdmin		bool		`jsonapi:"attr,is_admin"`
}



