package setucore

import (
	"OPQBot-QQGroupManager/Core"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Provider interface {
	InitProvider(l *logrus.Entry, bot *Core.Bot, db *gorm.DB)
	SearchPic(word string, r18 bool, num int) ([]Pic, error)
	SearchPicFromUser(word string, userId string, r18 bool, num int) ([]Pic, error)
}
type Pic struct {
	Id             int `gorm:"primaryKey"`
	Page           int `gorm:"primaryKey;default:0"`
	Title          string
	Author         string
	AuthorID       int
	OriginalPicUrl string
	Tag            string `gorm:"index;size:255"`
	R18            bool
	LastSendTime   int64 `gorm:"not null"`
}
