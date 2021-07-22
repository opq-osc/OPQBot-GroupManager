package setu

import (
	"gorm.io/gorm"
)

type DB struct {
	d *gorm.DB
}
type GroupConfig struct {
	GroupID int64 `gorm:"primaryKey"`
	R18     bool
	Enable  bool
}

func InitDB(db *gorm.DB) DB {
	//log.Println(db.AutoMigrate(&GroupConfig{},&setucore.Pic{}))
	return DB{d: db}
}
