package Local

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/GroupManager/Chat"

	"math/rand"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var log *logrus.Entry

type Core struct {
}
type LocalChat struct {
	gorm.Model
	GroupId  int64
	Question string
	Answer   string
	From     int64
}

func (c *Core) Init(l *logrus.Entry) error {
	log = l
	err := Config.DB.AutoMigrate(&LocalChat{})
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) GetAnswer(question string, GroupId, userId int64) (string, []byte) {
	var answers []LocalChat
	if err := Config.DB.Where("question LIKE ?", question).Find(&answers).Error; err != nil {
		log.Error(err)
		return "", nil
	}
	if len(answers) == 0 {
		return "", nil
	}
	return answers[rand.Intn(len(answers))].Answer, nil

}

func (c *Core) AddAnswer(question, answer string, GroupId, userId int64) error {
	var tmp = LocalChat{
		GroupId:  GroupId,
		Question: question,
		Answer:   answer,
		From:     userId,
	}
	if err := Config.DB.Create(&tmp).Error; err != nil {
		return err
	}
	return nil
}

func (c *Core) SetReplace(regexp string, target string) error {
	return nil
}

func init() {
	err := Chat.Register("local", &Core{})
	if err != nil {
		panic(err)
	}
}
