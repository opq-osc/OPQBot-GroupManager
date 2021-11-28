package setucore

import (
	"OPQBot-QQGroupManager/Core"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type Provider interface {
	InitProvider(l *logrus.Entry, bot *Core.Bot)
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

var (
	log *logrus.Entry
	db  *gorm.DB
)

func StartSetuCore(_db *gorm.DB, _log *logrus.Entry) {
	db = _db
	log = _log
	db.AutoMigrate(&Pic{})
}
func RegisterProvider(p Provider, log *logrus.Entry, bot *Core.Bot) {
	p.InitProvider(log, bot)
}
func SearchPicFromDB(word string, r18 bool, num int) (pics []Pic, e error) {
	if word == "" {
		e = db.Where("r18 = ? AND last_send_time < ?", r18, time.Now().Unix()-1800).Limit(num).Order("last_send_time asc").Find(&pics).Error
		return
	}
	e = db.Where("tag LIKE ? AND r18 = ? AND last_send_time < ?", "%"+word+"%", r18, time.Now().Unix()-1800).Limit(num).Order("last_send_time asc").Find(&pics).Error
	return
}
func SearchUserPicFromDB(word, userId string, r18 bool, num int) (pics []Pic, e error) {
	if userId != "" {
		e = db.Where("r18 = ? AND last_send_time < ? AND author_id = ?", r18, time.Now().Unix()-1800, userId).Limit(num).Order("last_send_time asc").Find(&pics).Error
		return
	}
	e = db.Where("author LIKE ? AND r18 = ? AND last_send_time < ?", "%"+word+"%", r18, time.Now().Unix()-1800).Limit(num).Order("last_send_time asc").Find(&pics).Error
	return
}
func AddPicToDB(pic Pic) error {
	var num int64
	db.Model(&pic).Where("id = ? AND page = ?", pic.Id, pic.Page).Count(&num)
	if num > 0 {
		return errors.New("该图片在数据库中已存在！")
	}
	return db.Create(&pic).Error
}
func PicInDB(picUrl string) bool {
	var num int64
	db.Model(&Pic{}).Where("original_pic_url = ?", picUrl).Count(&num)
	if num > 0 {
		return true
	}
	return false
}
func SetPicSendTime(pics []Pic) {
	for _, v := range pics {
		db.Model(&Pic{}).Where("id = ? AND page = ?", v.Id, v.Page).Updates(&Pic{LastSendTime: time.Now().Unix()})
	}

}
