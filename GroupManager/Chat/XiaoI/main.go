package XiaoI

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/GroupManager/Chat"
	"OPQBot-QQGroupManager/utils"
	"crypto/sha1"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
)

var log *logrus.Entry

func Sha1(data string) string {
	t := sha1.New()
	_, err := io.WriteString(t, data)
	if err != nil {
		log.Error(err)
		return ""
	}
	return fmt.Sprintf("%X", t.Sum(nil))
}

type Core struct {
}

func (c *Core) GetAnswer(question string, GroupId, userId int64) string {
	return ""
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
	key := Config.CoreConfig.ChatKey.XiaoI.Key
	secret := Config.CoreConfig.ChatKey.XiaoI.Secret
	Config.Lock.RUnlock()
	t1 := Sha1(key + ":xiaoi.com:" + secret)
	t2 := Sha1("POST:/ask.do")
	nonce := utils.RandomString(40)
	t3 := Sha1(t1 + ":" + nonce + ":" + t2)
	xAuth := fmt.Sprintf("app_key=\"%s\", nonce=\"%s\", signature=\"%s\"", key, nonce, t3)
	log.Info(xAuth)
	return nil
}
func init() {
	err := Chat.Register("XiaoI", &Core{})
	if err != nil {
		panic(err)
	}
}
