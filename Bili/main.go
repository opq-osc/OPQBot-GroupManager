package bili

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/Core"
	"encoding/base64"
	"fmt"
	"github.com/gaoyanpao/biliLiveHelper"
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/requests"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type Module struct {
}

var log *logrus.Entry

func (m *Module) ModuleInfo() Core.ModuleInfo {
	return Core.ModuleInfo{
		Name:        "Bili订阅姬",
		Author:      "enjoy",
		Description: "订阅bilibili番剧和UP",
		Version:     0,
	}
}
func (m *Module) ModuleInit(b *Core.Bot, l *logrus.Entry) error {
	log = l
	bi := NewManager()
	live := NewLiveManager()
	b.BotCronManager.AddJob(-1, "Bili", "*/5 * * * *", func() {
		update, fanju := bi.ScanUpdate()
		for _, v := range update {
			upName, gs, userId := bi.GetUpGroupsByMid(v.Mid)
			for _, g := range gs {
				if v1, ok := Config.CoreConfig.GroupConfig[g]; ok {
					if !v1.Bili {
						break
					}
				}
				res, _ := requests.Get(v.Pic)
				if userId == 0 {
					b.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  g,
						Content: OPQBot.SendTypePicMsgByBase64Content{
							Content: fmt.Sprintf("不知道是谁订阅的UP主%s更新了\n%s\n%s", upName, v.Title, v.Description),
							Base64:  base64.StdEncoding.EncodeToString(res.Content()),
							Flash:   false,
						},
					})
				} else {
					b.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  g,
						Content: OPQBot.SendTypePicMsgByBase64Content{
							Content: OPQBot.MacroAt([]int64{userId}) + fmt.Sprintf("您订阅的UP主%s更新了\n%s\n%s", upName, v.Title, v.Description),
							Base64:  base64.StdEncoding.EncodeToString(res.Content()),
							Flash:   false,
						},
					})
				}

			}
		}
		for _, v := range fanju {
			title, gs, userId := bi.GetFanjuGroupsByMid(v.Result.Media.MediaID)
			for _, g := range gs {
				if v1, ok := Config.CoreConfig.GroupConfig[g]; ok {
					if !v1.Bili {
						break
					}
				}
				res, _ := requests.Get(v.Result.Media.Cover)
				if userId == 0 {
					b.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  g,
						Content: OPQBot.SendTypePicMsgByBase64Content{
							Content: fmt.Sprintf("不知道是谁订阅的番剧%s更新了\n%s", title, v.Result.Media.NewEp.IndexShow),
							Base64:  base64.StdEncoding.EncodeToString(res.Content()),
							Flash:   false,
						},
					})
				} else {
					b.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  g,
						Content: OPQBot.SendTypePicMsgByBase64Content{
							Content: OPQBot.MacroAt([]int64{userId}) + fmt.Sprintf("您订阅的番剧%s更新了\n%s", title, v.Result.Media.NewEp.IndexShow),
							Base64:  base64.StdEncoding.EncodeToString(res.Content()),
							Flash:   false,
						},
					})
				}

			}
		}
	})
	err := b.AddEvent(OPQBot.EventNameOnGroupMessage, func(botQQ int64, packet *OPQBot.GroupMsgPack) {
		if packet.FromUserID == botQQ {
			return
		}
		Config.Lock.RLock()
		var c Config.GroupConfig
		if v, ok := Config.CoreConfig.GroupConfig[packet.FromGroupID]; ok {
			c = v
		} else {
			c = Config.CoreConfig.DefaultGroupConfig
		}
		Config.Lock.RUnlock()
		if !c.Enable {
			return
		}
		cm := strings.Split(packet.Content, " ")
		s := b.Session.SessionStart(packet.FromUserID)
		if packet.Content == "退出订阅" {
			err := s.Delete("biliUps")
			if err != nil {
				log.Println(err)
			}
			b.SendGroupTextMsg(packet.FromGroupID, "已经退出订阅")
			return
		}
		if packet.Content == "本群番剧" {
			if !c.Bili {
				return
			}
			ups := "本群订阅番剧\n"

			if len(c.Fanjus) == 0 {
				b.SendGroupTextMsg(packet.FromGroupID, "本群没有订阅番剧")
				return
			}
			for mid, v1 := range c.Fanjus {
				ups += fmt.Sprintf("%d - %s-订阅用户为：%d \n", mid, v1.Title, v1.UserId)
			}
			b.SendGroupTextMsg(packet.FromGroupID, ups)
		}
		if packet.Content == "本群up" {
			if !c.Bili {
				return
			}
			ups := "本群订阅UPs\n"

			if len(c.BiliUps) == 0 {
				b.SendGroupTextMsg(packet.FromGroupID, "本群没有订阅UP主")
				return
			}
			for mid, v1 := range c.BiliUps {
				ups += fmt.Sprintf("%d - %s -订阅者:%d\n", mid, v1.Name, v1.UserId)
			}
			b.SendGroupTextMsg(packet.FromGroupID, ups)
			return
		}
		if v, err := s.Get("biliUps"); err == nil {
			id, err := strconv.Atoi(packet.Content)
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, "序号错误, 输入“退出订阅”退出")
				return
			}
			if v1, ok := v.(map[int]int64); ok {
				if v2, ok := v1[id]; ok {
					u, err := bi.SubscribeUpByMid(packet.FromGroupID, v2, packet.FromUserID)
					if err != nil {
						b.SendGroupTextMsg(packet.FromGroupID, err.Error())
						err = s.Delete("biliUps")
						if err != nil {
							log.Println(err)
						}
						return
					}
					r, _ := requests.Get(u.Data.Card.Face)
					b.SendGroupPicMsg(packet.FromGroupID, "成功订阅UP主"+u.Data.Card.Name, r.Content())
					err = s.Delete("biliUps")
					if err != nil {
						log.Println(err)
					}
					return
				} else {
					b.SendGroupTextMsg(packet.FromGroupID, "序号不存在")
					return
				}

			} else {
				b.SendGroupTextMsg(packet.FromGroupID, "内部错误")
				err := s.Delete("biliUps")
				if err != nil {
					log.Println(err)
				}
				return
			}
			return
		}
		if len(cm) == 2 && cm[0] == "订阅up" {
			if !c.Bili {
				return
			}
			mid, err := strconv.ParseInt(cm[1], 10, 64)
			if err != nil {
				result, err := bi.SearchUp(cm[1])
				//u, err := bi.SubscribeUpByKeyword(packet.FromGroupID, cm[1])

				if err != nil {
					b.SendGroupTextMsg(packet.FromGroupID, err.Error())
					return
				}
				var (
					resultStr []string
					r         = map[int]int64{}
				)
				i := 0
				for _, v := range result.Data.Result {
					if v.IsUpuser == 1 {
						resultStr = append(resultStr, fmt.Sprintf("[%d] %s(lv.%d) 粉丝数:%d", i+1, v.Uname, v.Level, v.Fans))
						r[i+1] = v.Mid
						i++
						if len(r) >= 6 {
							break
						}
					}
				}
				if len(r) == 0 {
					b.SendGroupTextMsg(packet.FromGroupID, "没有找到UP哟~")
					return
				}
				err = s.Set("biliUps", r)
				if err != nil {
					b.SendGroupTextMsg(packet.FromGroupID, err.Error())
					return
				}
				b.SendGroupTextMsg(packet.FromGroupID, fmt.Sprintf("====输入序号选择UP====\n%s", strings.Join(resultStr, "\n")))
				return
			}
			u, err := bi.SubscribeUpByMid(packet.FromGroupID, mid, packet.FromUserID)
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			r, _ := requests.Get(u.Data.Card.Face)
			b.SendGroupPicMsg(packet.FromGroupID, "成功订阅UP主"+u.Data.Card.Name, r.Content())
			return
		}
		if len(cm) == 2 && cm[0] == "取消订阅up" {
			if !c.Bili {
				return
			}
			mid, err := strconv.ParseInt(cm[1], 10, 64)
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, "只能使用Mid取消订阅欧~")
				return
			}
			err = bi.UnSubscribeUp(packet.FromGroupID, mid)
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			b.SendGroupTextMsg(packet.FromGroupID, "成功取消订阅UP主")
			return
		}
		if len(cm) == 2 && cm[0] == "订阅番剧" {
			if !c.Bili {
				return
			}
			mid, err := strconv.ParseInt(cm[1], 10, 64)
			if err != nil {
				u, err := bi.SubscribeFanjuByKeyword(packet.FromGroupID, cm[1], packet.FromUserID)
				if err != nil {
					b.SendGroupTextMsg(packet.FromGroupID, err.Error())
					return
				}
				r, _ := requests.Get(u.Result.Media.Cover)
				b.SendGroupPicMsg(packet.FromGroupID, "成功订阅番剧"+u.Result.Media.Title, r.Content())
				return
			}
			u, err := bi.SubscribeFanjuByMid(packet.FromGroupID, mid, packet.FromUserID)
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			r, _ := requests.Get(u.Result.Media.Cover)
			b.SendGroupPicMsg(packet.FromGroupID, "成功订阅番剧"+u.Result.Media.Title, r.Content())
			return
		}
		if len(cm) == 2 && cm[0] == "取消订阅番剧" {
			if !c.Bili {
				return
			}
			mid, err := strconv.ParseInt(cm[1], 10, 64)
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, "只能使用Mid取消订阅欧~")
				return
			}
			err = bi.UnSubscribeFanju(packet.FromGroupID, mid)
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			b.SendGroupTextMsg(packet.FromGroupID, "成功取消订阅番剧")
			return
		}
		if len(cm) == 2 && cm[0] == "stopBili" {
			s, err := strconv.Atoi(cm[1])
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			err = live.RemoveClient(s)
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			b.SendGroupTextMsg(packet.FromGroupID, "已经断开连接了")
			return
		}
		if len(cm) == 2 && cm[0] == "biliLive" {
			if !Config.CoreConfig.BiliLive {
				b.SendGroupTextMsg(packet.FromGroupID, "该功能没有启动")
				return
			}
			s, err := strconv.Atoi(cm[1])
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			c, err := live.AddClient(s)
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			c.OnGift = func(ctx *biliLiveHelper.Context) {
				data := ctx.Msg
				r, _ := requests.Get(data.Get("data").Get("face").MustString())
				b.SendGroupPicMsg(packet.FromGroupID, fmt.Sprintf("%s%s%s", data.Get("data").Get("uname").MustString(), data.Get("data").Get("action").MustString(), data.Get("data").Get("giftName").MustString()), r.Content())
			}
			go c.Start()
			info := c.GetRoomInfo()
			b.SendGroupTextMsg(packet.FromGroupID, fmt.Sprintf("房间: %s[%d]\n关注: %d\n人气: %d\n直播状态: %v",
				info.Title,
				info.RoomID,
				info.Attention,
				info.Online,
				GetLiveStatusString(info.LiveStatus)))
			return
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func init() {
	Core.RegisterModule(&Module{})
}
