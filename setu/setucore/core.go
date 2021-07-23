package setucore

import (
	"OPQBot-QQGroupManager/Core"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Provider interface {
	InitProvider(l *logrus.Entry, bot *Core.Bot, db *gorm.DB)
	SearchPic(word string, r18 bool, num int) ([]Pic, error)
}
type Pic struct {
	Id             int `gorm:"primaryKey"`
	Title          string
	Author         string
	AuthorID       int
	OriginalPicUrl string
	Tag            string `gorm:"index"`
	R18            bool
	LastSendTime   int64 `gorm:"not null"`
}
