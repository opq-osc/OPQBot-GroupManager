package setucore

import (
	"OPQBot-QQGroupManager/Core"
	"github.com/sirupsen/logrus"
)

type Provider interface {
	InitProvider(l *logrus.Entry, bot *Core.Bot)
	SearchPic(word string, r18 bool) ([]Pic, error)
}
type Pic struct {
	Id             int `gorm:"primaryKey"`
	Title          string
	Author         string
	AuthorID       int
	OriginalPicUrl string
	Tag            string
	R18            bool
}
