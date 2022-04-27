package Zhai2

import (
	"OPQBot-QQGroupManager/GroupManager/Chat"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/sirupsen/logrus"
)

// ? 语料仓库未公开

var log *logrus.Entry

// https://raw.githubusercontent.com/opq-osc/awesome-2-methodology/main/resources/release/index.json
type Data struct {
	OnTrigger struct {
		BotName struct {
			Response []struct {
				Value string `json:"value"`
				Image string `json:"image,omitempty"`
			} `json:"response"`
		} `json:"botName"`
		Keyword []struct {
			Trigger  []string `json:"trigger"`
			Response []struct {
				Value string `json:"value,omitempty"`
				Image string `json:"image,omitempty"`
			} `json:"response"`
		} `json:"keyword"`
		Poke struct {
			Response []struct {
				Value string `json:"value"`
				Image string `json:"image"`
			} `json:"response"`
		} `json:"poke"`
	} `json:"onTrigger"`
}

var data Data

type Core struct {
}

func (c *Core) Init(l *logrus.Entry) error {
	log = l
	tmp, err := ioutil.ReadFile("./chat.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmp, &data)
	if err != nil {
		return err
	}
	return nil
}
func (c *Core) GetAnswer(question string, GroupId, userId int64) (string, []byte) {
	if question == "米娅" {
		a := data.OnTrigger.BotName.Response[rand.Intn(len(data.OnTrigger.BotName.Response))]
		var pic []byte
		if a.Image != "" {
			pic, _ = ioutil.ReadFile(a.Image)
			//pic = []byte(a.Image)
		}
		return strings.ReplaceAll(strings.ReplaceAll(a.Value, "{{bot.name}}", "米娅"), "{{user.nickname}}", "[YOU]"), pic
	}
	for _, v := range data.OnTrigger.Keyword {
		for _, v1 := range v.Trigger {
			if strings.Contains(question, strings.ReplaceAll(v1, "{{bot.name}}", "米娅")) {
				a := v.Response[rand.Intn(len(v.Response))]
				var pic []byte
				if a.Image != "" {
					pic, _ = ioutil.ReadFile(a.Image)
					//pic = []byte(a.Image)
				}
				return strings.ReplaceAll(strings.ReplaceAll(a.Value, "{{bot.name}}", "米娅"), "{{user.nickname}}", "[YOU]"), pic
			}
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
	err := Chat.Register("二次元高浓度", &Core{})
	if err != nil {
		panic(err)
	}
}
