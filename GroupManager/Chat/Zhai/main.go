package Zhai

import (
	"OPQBot-QQGroupManager/GroupManager/Chat"
	"math/rand"
	"strings"

	"github.com/mcoo/requests"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

type Core struct {
	Data map[string][]string
}

func (c *Core) Init(l *logrus.Entry) error {
	log = l
	r, err := requests.Get("https://cdn.jsdelivr.net/gh/Kyomotoi/AnimeThesaurus@main/data.json")
	if err != nil {
		return err
	}
	err = r.Json(&c.Data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) GetAnswer(question string, GroupId, userId int64) (string, []byte) {
	for k, v := range c.Data {
		if strings.Contains(question, k) {
			if len(v) == 0 {
				return "", nil
			}
			return v[rand.Intn(len(v))], nil
		}
	}
	return "", nil
}

func (c *Core) AddAnswer(question, answer string, GroupId, userId int64) error {
	return nil
}

func (c *Core) SetReplace(regexp string, target string) error {
	return nil
}

func init() {
	err := Chat.Register("二次元", &Core{})
	if err != nil {
		panic(err)
	}
}
