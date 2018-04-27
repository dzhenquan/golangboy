package model

import (
	"github.com/jinzhu/gorm"
	"gin-blog/config"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var DB *gorm.DB

// Base Model
type BaseModel struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", config.DBConfig.URL)

	if err == nil {
		DB = db
		db.SingularTable(true)
		db.AutoMigrate(&User{}, &Article{}, &Category{}, &Page{}, &Link{})
		return db, err
	}
	return nil, err
}