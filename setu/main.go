package setu

import (
	"OPQBot-QQGroupManager/Core"
	"OPQBot-QQGroupManager/setu/pixiv"
	"OPQBot-QQGroupManager/setu/setucore"
	"fmt"
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/requests"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

type Module struct {
}

var log *logrus.Entry

func (m *Module) ModuleInfo() Core.ModuleInfo {
	return Core.ModuleInfo{
		Name:        "Setu姬",
		Author:      "enjoy",
		Description: "思路来源于https://github.com/opq-osc/OPQ-SetuBot 天乐giegie的setu机器人",
		Version:     0,
	}
}
func (m *Module) ModuleInit(b *Core.Bot, l *logrus.Entry) error {
	log = l
	InitDB(b.DB)
	px := &pixiv.Provider{}
	RegisterProvider(px, b)
	err := b.AddEvent(OPQBot.EventNameOnGroupMessage, func(qq int64, packet *OPQBot.GroupMsgPack) {
		if packet.FromUserID != b.QQ {
			cm := strings.Split(packet.Content, " ")
			if len(cm) == 2 && cm[0] == "搜图测试" {
				pics, err := px.SearchPic(cm[1], false)
				rand.Seed(time.Now().UnixNano())
				num := rand.Intn(len(pics))
				res, err := requests.Get(strings.ReplaceAll(pics[num].OriginalPicUrl, "i.pximg.net", "i.pixiv.cat"))
				if err != nil {
					log.Error(err)
					return
				}
				b.SendGroupPicMsg(packet.FromGroupID, fmt.Sprintf("标题:%s", pics[num].Title), res.Content())
				if err != nil {
					log.Error(err)
				}

			}
		}
	})
	if err != nil {
		log.Error(err)
	}
	return nil
}
func init() {
	Core.RegisterModule(&Module{})
}
func RegisterProvider(p setucore.Provider, bot *Core.Bot) {
	p.InitProvider(log.WithField("SetuProvider", "Pixiv Core"), bot)
}
