package Moli

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/GroupManager/Chat"
	"errors"
	"github.com/mcoo/requests"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

type Core struct {
	replace map[string]string
	key     string
	secret  string
}

func (c *Core) GetAnswer(question string, GroupId, userId int64) string {
	res, err := requests.Get("http://i.itpk.cn/api.php?question=" + question + "&limit=8&api_key=" + c.key + "&api_secret=" + c.secret)
	if err != nil {
		log.Error(err)
		return ""
	}
	return res.Text()
}

func (c *Core) AddAnswer(question, answer string, GroupId, userId int64) error {
	return nil
}

func (c *Core) SetReplace(regexp string, target string) error {
	return nil
}

func (c *Core) Init(l *logrus.Entry) error {
	log = l
	Config.Lock.RLock()
	c.key = Config.CoreConfig.ChatKey.Moli.Key
	c.secret = Config.CoreConfig.ChatKey.Moli.Secret
	Config.Lock.RUnlock()
	c.replace = map[string]string{"[cqname]": "[YOU]", "[name]": "我"}
	if c.key == "" || c.secret == "" {
		return errors.New("key和密匙没有填写")
	}
	return nil
}
func init() {
	err := Chat.Register("茉莉", &Core{})
	if err != nil {
		panic(err)
	}
}
