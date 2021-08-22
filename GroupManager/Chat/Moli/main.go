package Moli

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/GroupManager/Chat"
	"errors"
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/requests"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var log *logrus.Entry

type Core struct {
	replace map[string]string
	key     string
	secret  string
}

func (c *Core) GetAnswer(question string, GroupId, userId int64) string {
	res, err := requests.Get("http://i.itpk.cn/api.php?question=" + OPQBot.DecodeFaceFromSentences(question, "%s") + "&limit=8&api_key=" + c.key + "&api_secret=" + c.secret)
	if err != nil {
		log.Error(err)
		return ""
	}
	if i, _ := regexp.MatchString(`http://`, res.Text()); i {
		return ""
	}
	return strings.ReplaceAll(res.Text(), "茉莉", "米娅")
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
	c.replace = map[string]string{"\\[cqname\\]": "\\[YOU\\]", "\\[name\\]": "米娅"}
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
