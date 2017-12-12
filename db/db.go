package db

import (
	"github.com/jinzhu/gorm"
)


var DatabaseConnectInstance *Db

type Db struct {
	db	*gorm.DB
}

func (db Db) GetDbConnect() *gorm.DB {
	return db.db
}

func NewDbConnect(connectionString string) *Db {
	db, _ := gorm.Open("mysql", connectionString)

	newDb := Db{db: db}
	DatabaseConnectInstance = &newDb

	return &newDb
}

func GetDatabaseConnectInstance() *Db {
	if DatabaseConnectInstance == nil {
		return NewDbConnect("")
	}

	return DatabaseConnectInstance

}