package setu

import (
	"OPQBot-QQGroupManager/Core"
	"OPQBot-QQGroupManager/setu/pixiv"
	"OPQBot-QQGroupManager/setu/setucore"
	"encoding/base64"
	"fmt"
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/requests"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Module struct {
}

var (
	log *logrus.Entry
)

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
	RegisterProvider(px, b, b.DB)
	err := b.AddEvent(OPQBot.EventNameOnGroupMessage, func(qq int64, packet *OPQBot.GroupMsgPack) {
		if packet.FromUserID != b.QQ {
			//cm := strings.Split(packet.Content, " ")
			cm := strings.SplitN(packet.Content, " ", 2)
			command, _ := regexp.Compile("搜([0-9]+)张图")
			if len(cm) >= 2 && command.MatchString(cm[0]) {
				getNum := 1
				tmp := command.FindStringSubmatch(cm[0])
				if len(tmp) == 2 {
					if count, err := strconv.Atoi(tmp[1]); err == nil && count > 1 {
						getNum = count
						if count > 10 {
							getNum = 10
						}
					}
				}

				pics, err := px.SearchPic(cm[1], false, getNum)
				if err != nil {
					log.Error(err)
					return
				}
				if len(pics) == 0 {
					b.SendGroupTextMsg(packet.FromGroupID, "没有找到有关"+cm[1]+"的图片")
					return
				}

				for _, v := range pics {
					res, err := requests.Get(strings.ReplaceAll(v.OriginalPicUrl, "i.pximg.net", "i-cf.pximg.net"), requests.Header{"referer": "https://www.pixiv.net/"})
					if err != nil {
						log.Error(err)
						return
					}
					b.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  packet.FromGroupID,
						Content: OPQBot.SendTypePicMsgByBase64Content{
							Content: fmt.Sprintf("标题:%s\n作者:%s\n作品链接:%s\n原图链接:%s\n30s自动撤回\n%s", v.Title, v.Author, "www.pixiv.net/artworks/"+strconv.Itoa(v.Id), v.OriginalPicUrl, OPQBot.MacroId()),
							Base64:  base64.StdEncoding.EncodeToString(res.Content()),
							Flash:   false,
						},
						CallbackFunc: func(Code int, Info string, record OPQBot.MyRecord) {
							time.Sleep(30 * time.Second)
							err := b.ReCallMsg(packet.FromGroupID, record.MsgRandom, record.MsgSeq)
							if err != nil {
								log.Warn(err)
							}
						},
					})
					if err != nil {
						log.Error(err)
					}
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
func RegisterProvider(p setucore.Provider, bot *Core.Bot, db *gorm.DB) {
	p.InitProvider(log.WithField("SetuProvider", "Pixiv Core"), bot, db)
}
