package GroupManager

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/Core"
	"OPQBot-QQGroupManager/GroupManager/Chat"
	_ "OPQBot-QQGroupManager/GroupManager/Chat/Local"
	//_ "OPQBot-QQGroupManager/GroupManager/Chat/XiaoI"
	_ "OPQBot-QQGroupManager/GroupManager/Chat/Moli"
	_ "OPQBot-QQGroupManager/GroupManager/Chat/Zhai"
	"OPQBot-QQGroupManager/GroupManager/QunInfo"
	"OPQBot-QQGroupManager/draw"
	"OPQBot-QQGroupManager/utils"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/mcoo/OPQBot"
	"github.com/sirupsen/logrus"
	"io/fs"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed Web/dist/spa
var staticFs embed.FS

type WebResult struct {
	Code int         `json:"code"`
	Info string      `json:"info"`
	Data interface{} `json:"data"`
}
type AtMsg struct {
	Content string `json:"Content"`
	UserExt []struct {
		QQNick string `json:"QQNick"`
		QQUID  int64  `json:"QQUid"`
	} `json:"UserExt"`
	UserID []int64 `json:"UserID"`
}

var (
	App   = iris.New()
	Start = make(chan struct{})
	sess  *sessions.Sessions
)

type Module struct {
}

var log *logrus.Entry

func (m *Module) ModuleInfo() Core.ModuleInfo {
	return Core.ModuleInfo{
		Name:        "群管理插件",
		Author:      "enjoy",
		Description: "",
		Version:     0,
	}
}
func (m *Module) ModuleInit(b *Core.Bot, l *logrus.Entry) error {
	log = l
	VerifyNum := map[string]*struct {
		Status bool
		Code   string
	}{}
	VerifyLock := sync.Mutex{}
	qun := QunInfo.NewQun(b)
	_, err := b.AddEvent(OPQBot.EventNameOnGroupJoin, func(botQQ int64, packet *OPQBot.GroupJoinPack) {
		Config.Lock.RLock()
		var c Config.GroupConfig
		if v, ok := Config.CoreConfig.GroupConfig[packet.EventMsg.FromUin]; ok {
			c = v
		} else {
			c = Config.CoreConfig.DefaultGroupConfig
		}
		Config.Lock.RUnlock()
		if !c.Enable {
			return
		}
		switch c.JoinVerifyType {
		case 1: // 图片验证码
			picB, n, err := draw.Draw6Number()
			if err != nil {
				log.Println(err)
				return
			}
			b.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeGroup,
				ToUserUid:  packet.EventMsg.FromUin,
				Content: OPQBot.SendTypePicMsgByBase64Content{
					Content: OPQBot.MacroAt([]int64{packet.EventData.UserID}) + "请在" + strconv.Itoa(c.JoinVerifyTime) + "s内输入上方图片验证码！否则会被移出群,若看不清楚可以输入 刷新验证码\n" + OPQBot.MacroId(),
					Base64:  base64.StdEncoding.EncodeToString(picB),
					Flash:   false,
				},
				CallbackFunc: func(Code int, Info string, record OPQBot.MyRecord) {
					if record.MsgSeq == 0 {
						log.Println("验证码信息没有发送成功！")
					} else {
						VerifyLock.Lock()
						VerifyNum[strconv.FormatInt(packet.EventData.UserID, 10)+"|"+strconv.FormatInt(packet.EventMsg.FromUin, 10)] = &struct {
							Status bool
							Code   string
						}{Status: false, Code: n}
						VerifyLock.Unlock()
						time.Sleep(time.Duration(c.JoinVerifyTime) * time.Second)
						VerifyLock.Lock()
						log.Println(VerifyNum)
						if v, ok := VerifyNum[strconv.FormatInt(packet.EventData.UserID, 10)+"|"+strconv.FormatInt(packet.EventMsg.FromUin, 10)]; ok {
							if !v.Status {
								b.Send(OPQBot.SendMsgPack{
									SendToType: OPQBot.SendToTypeGroup,
									ToUserUid:  packet.EventMsg.FromUin,
									Content: OPQBot.SendTypeTextMsgContent{
										Content: OPQBot.MacroAt([]int64{packet.EventData.UserID}) + "验证超时,再见!",
									},
									CallbackFunc: func(Code int, Info string, record OPQBot.MyRecord) {
										b.KickGroupMember(packet.EventMsg.FromUin, packet.EventData.UserID)
									},
								})
							} else {
								delete(VerifyNum, strconv.FormatInt(packet.EventData.UserID, 10)+"|"+strconv.FormatInt(packet.EventMsg.FromUin, 10))
							}
						}
						VerifyLock.Unlock()
					}
				},
			})

		default:
		}
		if c.Welcome != "" {
			b.SendGroupTextMsg(packet.EventMsg.FromUin, c.Welcome)
		}
		if c.JoinAutoShutUpTime != 0 {
			b.SetForbidden(1, c.JoinAutoShutUpTime, packet.EventMsg.FromUin, packet.EventData.UserID)
		}
	})
	// 接受处理解禁功能
	_, err = b.AddEvent(OPQBot.EventNameOnFriendMessage, func(qq int64, packet *OPQBot.FriendMsgPack) {
		if packet.FromUin != b.QQ {
			//log.Println(packet.Content)
			content := map[string]interface{}{}
			contentStr := ""
			json.Unmarshal([]byte(packet.Content), &content)
			if value, ok := content["Content"]; ok {
				contentStr, _ = value.(string)
			} else {
				contentStr = packet.Content
			}
			s := b.Session.SessionStart(packet.FromUin)
			if result, err := s.GetInt("c_result"); err == nil {
				IgroupId, err := s.Get("groupId")
				if err != nil {
					log.Error(err)
					return
				}
				groupId, _ := IgroupId.(int64)
				if contentStr == "取消" {
					s.Delete("groupId")
					s.Delete("c_result")
					b.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypePrivateChat,
						ToUserUid:  packet.FromUin,
						Content: OPQBot.SendTypeTextMsgContentPrivateChat{
							Content: "已经取消了",
							Group:   groupId,
						},
						CallbackFunc: nil,
					})
					//b.SendFriendTextMsg(packet.FromUin,"已经取消了！")
					return
				}
				if strconv.Itoa(result) == contentStr {

					err = packet.Bot.SetForbidden(1, 0, groupId, packet.FromUin)
					if err != nil {
						log.Error(err)
						return
					}
					b.SendGroupTextMsg(groupId, fmt.Sprintf("用户%d,解除禁言", packet.FromUin))
					//b.SendFriendTextMsg(packet.FromUin,"已经操作了！")
					b.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypePrivateChat,
						ToUserUid:  packet.FromUin,
						Content: OPQBot.SendTypeTextMsgContentPrivateChat{
							Content: "已经操作了",
							Group:   groupId,
						},
						CallbackFunc: nil,
					})
					s.Delete("groupId")
					s.Delete("c_result")
				} else {
					b.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypePrivateChat,
						ToUserUid:  packet.FromUin,
						Content: OPQBot.SendTypeTextMsgContentPrivateChat{
							Content: "答案错误！输入“取消”退出",
							Group:   groupId,
						},
						CallbackFunc: nil,
					})
					//b.SendFriendTextMsg(packet.FromUin,"答案错误！输入“取消”退出")
					return
				}
			}
			if find := strings.Contains(packet.Content, "解除禁言"); find {
				cm := strings.Split(contentStr, ",")
				if len(cm) != 2 {
					return
				}
				groupId, err := strconv.ParseInt(cm[1], 10, 64)

				rand.Seed(time.Now().Unix())
				a1 := rand.Intn(100)
				a2 := rand.Intn(100)
				err = s.Set("c_result", a1+a2)
				if err != nil {
					log.Error(err)
					return
				}
				err = s.Set("groupId", groupId)
				if err != nil {
					log.Error(err)
					return
				}
				b.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypePrivateChat,
					ToUserUid:  packet.FromUin,
					Content: OPQBot.SendTypeTextMsgContentPrivateChat{
						Content: fmt.Sprintf("你好,请先回答对问题才能解除禁言哟！问题:\n%s", base64.StdEncoding.EncodeToString([]byte(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d + %d = ?", a1, a2)))))),
						Group:   groupId,
					},
					CallbackFunc: nil,
				})
				//b.SendFriendTextMsg(packet.FromUin, fmt.Sprintf("你好,请先回答对问题才能解除禁言哟！问题:\n%s",base64.StdEncoding.EncodeToString([]byte(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d + %d = ?",a1,a2)))))))

			}
		}
	})
	if err != nil {
		log.Error(err)
	}
	chat := Chat.StartChatCore(log.WithField("Func", "Chat"))
	_, err = b.AddEvent(OPQBot.EventNameOnGroupMessage, func(botQQ int64, packet *OPQBot.GroupMsgPack) {
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

		if v, _ := regexp.MatchString(`^[0-9]{6}$`, packet.Content); v {
			VerifyLock.Lock()
			if v1, ok := VerifyNum[strconv.FormatInt(packet.FromUserID, 10)+"|"+strconv.FormatInt(packet.FromGroupID, 10)]; ok {
				if v1.Code == packet.Content {
					v1.Status = true
					b.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  packet.FromGroupID,
						Content: OPQBot.SendTypeTextMsgContent{
							Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "验证成功",
						},
					})
				}
			}
			VerifyLock.Unlock()
			return
		}
		if m, err := regexp.MatchString(c.MenuKeyWord, packet.Content); err != nil {
			log.Println(err)
			return
		} else if m {
			b.Send(OPQBot.SendMsgPack{
				SendToType:   OPQBot.SendToTypeGroup,
				ToUserUid:    packet.FromGroupID,
				Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + c.Menu},
				CallbackFunc: nil,
			})
			return
		}

		if packet.Content == "刷新验证码" {
			VerifyLock.Lock()
			if v1, ok := VerifyNum[strconv.FormatInt(packet.FromUserID, 10)+"|"+strconv.FormatInt(packet.FromGroupID, 10)]; ok {
				picB, n, err := draw.Draw6Number()
				if err != nil {
					log.Println(err)
					return
				}
				v1.Code = n
				VerifyNum[strconv.FormatInt(packet.FromUserID, 10)+"|"+strconv.FormatInt(packet.FromGroupID, 10)] = v1
				b.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content: OPQBot.SendTypePicMsgByBase64Content{
						Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "请在" + strconv.Itoa(c.JoinVerifyTime) + "s内输入上方图片验证码,全是数字哟！否则会被移出群,若看不清楚可以输入 刷新验证码\n" + OPQBot.MacroId(),
						Base64:  base64.StdEncoding.EncodeToString(picB),
						Flash:   false,
					},
				})
			}
			VerifyLock.Unlock()
			return
		}
		if packet.Content == "本群信息" {
			info, err := qun.GetGroupInfo(packet.FromGroupID, 0)
			if err != nil {
				log.Error(err)
				return
			}
			info2, err := qun.GetGroupMembersInfo(packet.FromGroupID, 0)
			if err != nil {
				log.Error(err)
				return
			}
			pic, err := draw.DrawGroupInfo(info, info2)
			if err != nil {
				log.Error(err)
				return
			}
			b.SendGroupPicMsg(packet.FromGroupID, "", pic)
			return
			//s := fmt.Sprintf(
			//	"本群[%d]%s人数%d\n昨日活跃数据:\n活跃人数:%d\n消息条数:%d\n加群%d人 退群%d人 申请入群%d\n最活跃的小可爱们",
			//	packet.FromGroupID,
			//	info.Data.GroupInfo.GroupName,
			//	info.Data.GroupInfo.GroupMember,
			//	info.Data.ActiveData.DataList[len(info.Data.ActiveData.DataList)-1].Number,
			//	info.Data.MsgInfo.DataList[len(info.Data.MsgInfo.DataList)-1].Number,
			//	info.Data.JoinData.DataList[len(info.Data.JoinData.DataList)-1].Number,
			//	info.Data.ExitData.DataList[len(info.Data.ExitData.DataList)-1].Number,
			//	info.Data.ApplyData.DataList[len(info.Data.ApplyData.DataList)-1].Number,
			//)
			//a := 0
			//for _, v := range info2.Data.SpeakRank {
			//	if a >= 5 {
			//		break
			//	}
			//
			//	s += fmt.Sprintf("\n%s 活跃度 %d 发言条数 %d", v.Nickname, v.Active, v.MsgCount)
			//	a += 1
			//}
			//b.SendGroupTextMsg(packet.FromGroupID, s)
			//return
		}
		if packet.Content == "签到" {

			if !c.SignIn {
				b.Send(OPQBot.SendMsgPack{
					SendToType:   OPQBot.SendToTypeGroup,
					ToUserUid:    packet.FromGroupID,
					Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "本群签到功能未开启!"},
					CallbackFunc: nil,
				})
				return
			}
			Config.Lock.Lock()
			if v, ok := Config.CoreConfig.UserData[packet.FromUserID]; ok {
				if v.LastSignDay == time.Now().Day() {
					b.Send(OPQBot.SendMsgPack{
						SendToType:   OPQBot.SendToTypeGroup,
						ToUserUid:    packet.FromGroupID,
						Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "今日已经签到过了,明日再来"},
						CallbackFunc: nil,
					})
				} else {
					v.Count += 1
					v.LastSignDay = time.Now().Day()
					Config.CoreConfig.UserData[packet.FromUserID] = v
					err := Config.Save()
					if err != nil {
						log.Println(err)
					}
					b.Send(OPQBot.SendMsgPack{
						SendToType:   OPQBot.SendToTypeGroup,
						ToUserUid:    packet.FromGroupID,
						Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "签到成功"},
						CallbackFunc: nil,
					})
				}
			} else {
				v.Count = 1
				v.LastSignDay = time.Now().Day()
				Config.CoreConfig.UserData[packet.FromUserID] = v
				err := Config.Save()
				if err != nil {
					log.Println(err)
				}
				b.Send(OPQBot.SendMsgPack{
					SendToType:   OPQBot.SendToTypeGroup,
					ToUserUid:    packet.FromGroupID,
					Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "签到成功"},
					CallbackFunc: nil,
				})
			}
			Config.Lock.Unlock()
			return
		}
		if packet.Content == "赞我" {
			if !c.Zan {
				b.Send(OPQBot.SendMsgPack{
					SendToType:   OPQBot.SendToTypeGroup,
					ToUserUid:    packet.FromGroupID,
					Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "本群名片赞功能未开启!"},
					CallbackFunc: nil,
				})
				return
			}
			Config.Lock.Lock()
			if v, ok := Config.CoreConfig.UserData[packet.FromUserID]; ok {
				if v.LastZanDay == time.Now().Day() {
					b.Send(OPQBot.SendMsgPack{
						SendToType:   OPQBot.SendToTypeGroup,
						ToUserUid:    packet.FromGroupID,
						Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "今日已经赞过了,明日再来"},
						CallbackFunc: nil,
					})
				} else {
					v.LastZanDay = time.Now().Day()
					Config.CoreConfig.UserData[packet.FromUserID] = v
					err := Config.Save()
					if err != nil {
						log.Println(err)
					}
					b.Send(OPQBot.SendMsgPack{
						SendToType:   OPQBot.SendToTypeGroup,
						ToUserUid:    packet.FromGroupID,
						Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "正在赞请稍后"},
						CallbackFunc: nil,
					})
				}
			} else {
				v.LastZanDay = time.Now().Day()
				Config.CoreConfig.UserData[packet.FromUserID] = v
				err := Config.Save()
				if err != nil {
					log.Println(err)
				}
				b.Send(OPQBot.SendMsgPack{
					SendToType:   OPQBot.SendToTypeGroup,
					ToUserUid:    packet.FromGroupID,
					Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "正在赞请稍后"},
					CallbackFunc: nil,
				})
			}
			Config.Lock.Unlock()
			return
		}
		if packet.Content == "积分" {
			Config.Lock.RLock()
			if v, ok := Config.CoreConfig.UserData[packet.FromUserID]; ok {
				b.Send(OPQBot.SendMsgPack{
					SendToType:   OPQBot.SendToTypeGroup,
					ToUserUid:    packet.FromGroupID,
					Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "你的积分为" + strconv.Itoa(v.Count)},
					CallbackFunc: nil,
				})
			} else {
				b.Send(OPQBot.SendMsgPack{
					SendToType:   OPQBot.SendToTypeGroup,
					ToUserUid:    packet.FromGroupID,
					Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "你的积分为0"},
					CallbackFunc: nil,
				})
			}
			Config.Lock.RUnlock()
			return
		}
		if a, _ := regexp.MatchString("聊天|无聊", packet.Content); a {
			Config.Lock.Lock()
			c.EnableChat = true
			Config.CoreConfig.GroupConfig[packet.FromGroupID] = c
			Config.Save()
			Config.Lock.Unlock()
			b.SendGroupTextMsg(packet.FromGroupID, "你可以at我和我聊天哟😉")
			return
		}
		if packet.Content == "当前聊天数据库" {
			if chat.SelectCore == "" {
				b.SendGroupTextMsg(packet.FromGroupID, "当前没有设置")
				return
			}
			b.SendGroupTextMsg(packet.FromGroupID, "设置聊天数据库为"+chat.SelectCore)
			return
		}
		cm := strings.Split(packet.Content, " ")
		if len(cm) == 2 && cm[0] == "拉黑" {
			Config.Lock.Lock()
			log.Println(cm)
			if Config.CoreConfig.SuperAdminUin == packet.FromUserID {
				if uin, err := strconv.ParseInt(cm[1], 10, 64); err == nil {
					Config.CoreConfig.BanQQ = append(Config.CoreConfig.BanQQ, uin)
					b.SendGroupTextMsg(packet.FromGroupID, "已经拉黑了")
				} else {
					log.Error(err)
					b.SendGroupTextMsg(packet.FromGroupID, "命令有问题")
				}

			} else {
				b.SendGroupTextMsg(packet.FromGroupID, "权限不足")
			}
			Config.Save()
			Config.Lock.Unlock()
			return
		}
		if len(cm) == 2 && cm[0] == "取消拉黑" {
			Config.Lock.Lock()
			if Config.CoreConfig.SuperAdminUin == packet.FromUserID {
				if uin, err := strconv.ParseInt(cm[1], 10, 64); err != nil {
					for i, v := range Config.CoreConfig.BanQQ {
						if v == uin {
							Config.CoreConfig.BanQQ = append(Config.CoreConfig.BanQQ[:i], Config.CoreConfig.BanQQ[i+1:]...)
						}
					}
					b.SendGroupTextMsg(packet.FromGroupID, "已经拉黑了")
				} else {
					b.SendGroupTextMsg(packet.FromGroupID, "命令有问题")
				}
			} else {
				b.SendGroupTextMsg(packet.FromGroupID, "权限不足")
			}
			Config.Save()
			Config.Lock.Unlock()
			return
		}
		if packet.Content == "黑名单" {
			Config.Lock.RLock()
			s := "黑名单列表"
			for _, v := range Config.CoreConfig.BanQQ {
				s += "\n" + strconv.FormatInt(v, 10)
			}
			b.SendGroupTextMsg(packet.FromGroupID, s)
			Config.Lock.RUnlock()
			return
		}
		if len(cm) >= 3 && cm[0] == "教你" {
			tmp := strings.SplitN(packet.Content, " ", 3)
			err := chat.Learn(tmp[1], tmp[2], packet.FromGroupID, packet.FromUserID)
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, "学习出现问题")
				return
			}
			b.SendGroupTextMsg(packet.FromGroupID, "已经学会了")
			return
		}
		if len(cm) == 2 && cm[0] == "设置聊天数据库" {
			err = chat.SetChatDB(cm[1])
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			b.SendGroupTextMsg(packet.FromGroupID, "设置聊天数据库为"+cm[1])
			return
		}
		if a, _ := regexp.MatchString("关闭|无路赛|吵|别说话|闭嘴", packet.Content); a {
			Config.Lock.Lock()
			if c.Enable == true {
				c.EnableChat = false
				b.SendGroupTextMsg(packet.FromGroupID, "关闭了聊天功能")
				Config.CoreConfig.GroupConfig[packet.FromGroupID] = c
				Config.Save()
			}
			Config.Lock.Unlock()
			return
		}
		if c.EnableChat && packet.MsgType == "AtMsg" {
			//var atInfo AtMsg
			//err := json.Unmarshal([]byte(packet.Content),&atInfo)
			//if err != nil{
			//	log.Error(err)
			//	return
			//}
			atInfo, err := OPQBot.ParserGroupAtMsg(*packet)
			if err != nil {
				log.Error(err)
				return
			}
			atInfo = atInfo.Clean()
			var atme = false
			for _, v := range atInfo.UserID {
				if v == b.QQ {
					atme = true
					break
				}
			}
			if !atme {
				return
			}
			for _, v := range atInfo.UserExt {
				atInfo.Content = strings.TrimSpace(strings.ReplaceAll(atInfo.Content, "@"+v.QQNick, ""))
			}

			answer, err := chat.GetAnswer(OPQBot.DecodeFaceFromSentences(atInfo.Content, "%s"), packet.FromGroupID, packet.FromUserID)
			if err != nil {
				log.Warn(err)
				return
			}
			b.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeGroup,
				ToUserUid:  packet.FromGroupID,
				Content: OPQBot.SendTypeReplyContent{
					ReplayInfo: struct {
						MsgSeq     int    `json:"MsgSeq"`
						MsgTime    int    `json:"MsgTime"`
						UserID     int64  `json:"UserID"`
						RawContent string `json:"RawContent"`
					}{MsgSeq: packet.MsgSeq,
						MsgTime:    packet.MsgTime,
						UserID:     packet.FromUserID,
						RawContent: atInfo.Content,
					},
					Content: strings.ReplaceAll(answer, "[YOU]", packet.FromNickName),
				},
				CallbackFunc: nil,
			})
			//b.SendGroupTextMsg(packet.FromGroupID, )
		}
	})
	if err != nil {
		return err
	}

	if Config.CoreConfig.OPQWebConfig.Enable {
		log.Println("启动Web 😊")
		Config.Lock.Lock()
		sess = sessions.New(sessions.Config{Cookie: "OPQWebSession"})
		if Config.CoreConfig.OPQWebConfig.CSRF == "" {
			Config.CoreConfig.OPQWebConfig.CSRF = utils.RandomString(32)
			err := Config.Save()
			if err != nil {
				log.Println(err)
			}
		}
		fads, _ := fs.Sub(staticFs, "Web/dist/spa")
		if Config.CoreConfig.ReverseProxy != "" {
			// target, err := url.Parse(Config.CoreConfig.ReverseProxy)
			if err != nil {
				panic(err)
			}
			App.Get("{root:path}", func(ctx iris.Context) {
				director := func(r *http.Request) {
					r.Host = Config.CoreConfig.ReverseProxy
					r.URL, _ = url.Parse(r.Host + "/" + ctx.Path())
				}
				p := &httputil.ReverseProxy{Director: director}
				p.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
			})
		} else {
			App.HandleDir("/", http.FS(fads))
		}

		// app.HandleDir("/", iris.Dir("./Web/dist/spa"))
		Config.Lock.Unlock()

		App.Use(beforeCsrf)
		App.Use(sess.Handler())
		App.WrapRouter(func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
			w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			path := r.URL.Path
			if r.Method == "OPTIONS" {
				w.Header().Add("Access-Control-Allow-Headers", "content-type")
				w.WriteHeader(200)
				return
			}
			if len(path) < 4 {
				if !pathIsFile(path) {
					r.URL.Path = "/"
				}
			} else {
				if r.URL.Path[0:4] != "/api" && r.URL.Path[0:4] != "/git" && r.URL.Path[0:4] != "/wor" {
					if !pathIsFile(path) {
						r.URL.Path = "/"
					}
				}
			}
			// log.Println(r.URL.Path)
			router.ServeHTTP(w, r)
		})
		App.Get("/api/chat", func(ctx iris.Context) {
			word := ctx.URLParamDefault("s", "")
			if word == "" {
				ctx.JSON(WebResult{
					Code: 0,
					Info: "success",
					Data: "你想要问米娅什么事情呢？？",
				})
				return
			}
			a, err := chat.GetAnswer(word, 0, 0)
			if err != nil {
				ctx.JSON(WebResult{
					Code: 1,
					Info: "error",
					Data: err.Error(),
				})
				return
			}
			ctx.JSON(WebResult{
				Code: 0,
				Info: "success",
				Data: a,
			})
			return
		})
		App.Get("/api/csrf", func(ctx iris.Context) {
			s := sess.Start(ctx)
			salt := int(time.Now().Unix())
			keyTmp := utils.Md5V(strconv.Itoa(salt + rand.Intn(100)))
			s.Set("OPQWebCSRF", keyTmp)
			ctx.SetCookieKV("OPQWebCSRF", keyTmp, iris.CookieHTTPOnly(false))
			_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: s.Get("username")})
		})
		App.Get("/api/status", func(ctx iris.Context) {
			s := sess.Start(ctx)
			salt := int(time.Now().Unix())
			keyTmp := utils.Md5V(strconv.Itoa(salt + rand.Intn(100)))
			s.Set("OPQWebCSRF", keyTmp)
			ctx.SetCookieKV("OPQWebCSRF", keyTmp, iris.CookieHTTPOnly(false))
			if s.GetBooleanDefault("auth", false) {
				_, _ = ctx.JSON(WebResult{Code: 1, Info: "已登录!", Data: s.Get("username")})
				return
			} else {
				_, _ = ctx.JSON(WebResult{Code: 0, Info: "未登录!", Data: nil})
				return
			}
		})
		App.Post("/api/login", func(ctx iris.Context) {
			username := ctx.FormValue("username")
			password := ctx.FormValue("password")
			Config.Lock.RLock()
			defer Config.Lock.RUnlock()
			if username == Config.CoreConfig.OPQWebConfig.Username && password == utils.Md5V(Config.CoreConfig.OPQWebConfig.Password) {
				s := sess.Start(ctx)
				s.Set("auth", true)
				_, _ = ctx.JSON(WebResult{Code: 1, Info: "登录成功", Data: nil})
				return
			} else {
				_, _ = ctx.JSON(WebResult{Code: 0, Info: "用户名密码错误!", Data: nil})
				return
			}

		})
		// job周期任务读取
		Config.Lock.RLock()
		for k, v := range Config.CoreConfig.GroupConfig {
			for k1, v2 := range v.Job {
				switch v2.Type {
				case 1:
					err = b.BotCronManager.AddJob(k, k1, v2.Cron, func() {
						log.Print("执行任务" + k1)
						if b.Announce(v2.Title, v2.Content, 0, 10, k) != nil {
							log.Print(err)
						}
					})
					if err != nil {
						log.Print("添加任务" + k1 + "出现错误" + err.Error())
					}
				case 2:
					err = b.BotCronManager.AddJob(k, k1, v2.Cron, func() {
						log.Print("执行任务" + k1)
						if b.SetForbidden(0, 1, k, 0) != nil {
							log.Print(err)
						}
					})
					if err != nil {
						log.Print("添加任务" + k1 + "出现错误" + err.Error())
					}
				case 3:
					err = b.BotCronManager.AddJob(k, k1, v2.Cron, func() {
						log.Print("执行任务" + k1)
						if b.SetForbidden(0, 0, k, 0) != nil {
							log.Print(err)
						}
					})
					if err != nil {
						log.Print("添加任务" + k1 + "出现错误" + err.Error())
					}
				case 4:
					err = b.BotCronManager.AddJob(k, k1, v2.Cron, func() {
						log.Print("执行任务" + k1)
						b.Send(OPQBot.SendMsgPack{
							SendToType: OPQBot.SendToTypeGroup,
							ToUserUid:  k,
							Content: OPQBot.SendTypeTextMsgContent{
								Content: v2.Content,
							},
						})
					})
					if err != nil {
						log.Print("添加任务" + k1 + "出现错误" + err.Error())
					}
				}
			}
		}
		Config.Lock.RUnlock()
		needAuth := App.Party("/api/admin", requireAuth)
		{
			rJob := needAuth.Party("/job")
			{
				rJob.Post("/add", func(ctx iris.Context) {
					ids := ctx.FormValue("id")
					id, err := strconv.ParseInt(ids, 10, 64)
					if err != nil {
						_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
						return
					}
					if id == -1 {
						_, _ = ctx.JSON(WebResult{Code: 0, Info: "默认群禁止添加周期任务", Data: nil})
						return
					}
					span := ctx.FormValue("span")
					jobName := ctx.FormValue("jobName")
					if jobName == "" {
						_, _ = ctx.JSON(WebResult{Code: 0, Info: "jobName为空", Data: nil})
						return
					}
					cronType, _ := strconv.Atoi(ctx.FormValue("type"))
					switch cronType {
					// 公告
					case 1:
						title := ctx.FormValue("title")
						content := ctx.FormValue("content")
						err = b.BotCronManager.AddJob(id, jobName, span, func() {
							log.Print("执行任务" + jobName)
							if b.Announce(title, content, 0, 10, id) != nil {
								log.Print(err)
							}
						})
						if err != nil {
							_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
							return
						} else {
							job := Config.Job{Type: cronType, Cron: span, Title: title, Content: content}
							Config.Lock.Lock()
							defer Config.Lock.Unlock()
							if v, ok := Config.CoreConfig.GroupConfig[id]; ok {
								v.Job[jobName] = job
								Config.CoreConfig.GroupConfig[id] = v
							} else {
								v = Config.CoreConfig.DefaultGroupConfig
								v.Job[jobName] = job
								Config.CoreConfig.GroupConfig[id] = v
							}
							Config.Save()
							_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: nil})
							return
						}
					// 全局禁言
					case 2:
						err = b.BotCronManager.AddJob(id, jobName, span, func() {
							log.Print("执行任务" + jobName)
							if b.SetForbidden(0, 1, id, 0) != nil {
								log.Print(err)
							}
						})
						if err != nil {
							_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
							return
						} else {
							job := Config.Job{Type: cronType, Cron: span}
							Config.Lock.Lock()
							defer Config.Lock.Unlock()
							if v, ok := Config.CoreConfig.GroupConfig[id]; ok {
								v.Job[jobName] = job
								Config.CoreConfig.GroupConfig[id] = v
							} else {
								v = Config.CoreConfig.DefaultGroupConfig
								v.Job[jobName] = job
								Config.CoreConfig.GroupConfig[id] = v
							}
							Config.Save()
							_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: nil})
							return
						}
					// 解除全局禁言
					case 3:
						err = b.BotCronManager.AddJob(id, jobName, span, func() {
							log.Print("执行任务" + jobName)
							if b.SetForbidden(0, 0, id, 0) != nil {
								log.Print(err)
							}
						})
						if err != nil {
							_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
							return
						} else {
							job := Config.Job{Type: cronType, Cron: span}
							Config.Lock.Lock()
							defer Config.Lock.Unlock()
							if v, ok := Config.CoreConfig.GroupConfig[id]; ok {
								v.Job[jobName] = job
								Config.CoreConfig.GroupConfig[id] = v
							} else {
								v = Config.CoreConfig.DefaultGroupConfig
								v.Job[jobName] = job
								Config.CoreConfig.GroupConfig[id] = v
							}
							Config.Save()
							_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: nil})
							return
						}
					case 4:
						// title := ctx.FormValue("title")
						content := ctx.FormValue("content")
						err = b.BotCronManager.AddJob(id, jobName, span, func() {
							log.Print("执行任务" + jobName)
							b.Send(OPQBot.SendMsgPack{
								SendToType: OPQBot.SendToTypeGroup,
								ToUserUid:  id,
								Content: OPQBot.SendTypeTextMsgContent{
									Content: content,
								},
							})
						})
						if err != nil {
							_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
							return
						} else {
							job := Config.Job{Type: cronType, Cron: span, Content: content}
							Config.Lock.Lock()
							defer Config.Lock.Unlock()
							if v, ok := Config.CoreConfig.GroupConfig[id]; ok {
								v.Job[jobName] = job
								Config.CoreConfig.GroupConfig[id] = v
							} else {
								v = Config.CoreConfig.DefaultGroupConfig
								v.Job[jobName] = job
								Config.CoreConfig.GroupConfig[id] = v
							}
							Config.Save()
							_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: nil})
							return
						}
					default:
						_, _ = ctx.JSON(WebResult{Code: 0, Info: "类型不存在", Data: nil})
						return
					}
				})
				rJob.Post("/del", func(ctx iris.Context) {
					ids := ctx.FormValue("id")
					id, err := strconv.ParseInt(ids, 10, 64)
					if err != nil {
						_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
						return
					}
					if id == -1 {
						_, _ = ctx.JSON(WebResult{Code: 0, Info: "默认群禁止删除周期任务", Data: nil})
						return
					}
					jobName := ctx.FormValue("jobName")
					if jobName == "" {
						_, _ = ctx.JSON(WebResult{Code: 0, Info: "jobName为空", Data: nil})
						return
					}
					err = b.BotCronManager.Remove(id, jobName)

					if err != nil {
						_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
						return
					}
					Config.Lock.Lock()
					defer Config.Lock.Unlock()
					if v, ok := Config.CoreConfig.GroupConfig[id]; ok {
						delete(v.Job, jobName)
						Config.CoreConfig.GroupConfig[id] = v
					} else {
						_, _ = ctx.JSON(WebResult{Code: 0, Info: "Group在配置文件中不存在！", Data: nil})
						return
					}
					Config.Save()
					_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: nil})
				})
			}
			needAuth.Post("/getGroupMember", func(ctx iris.Context) {
				ids := ctx.FormValue("id")
				id, err := strconv.ParseInt(ids, 10, 64)
				if err != nil {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
					return
				}
				if id == -1 {
					_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: []int{}})
					return
				}

				glist, err := b.GetGroupMemberList(id, 0)
				if err != nil {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
					return
				}
				result := glist
				for {
					if glist.LastUin == 0 {
						break
					}
					glist, err = b.GetGroupMemberList(id, glist.LastUin)
					if err != nil {
						_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
						return
					}
					result.MemberList = append(result.MemberList, glist.MemberList...)
					result.Count += glist.Count
					result.LastUin = glist.LastUin
				}
				_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: result})
				return
			})
			needAuth.Post("/setGroupConfig", func(ctx iris.Context) {
				ids := ctx.FormValue("id")
				enable := ctx.FormValue("enable")
				id, err := strconv.ParseInt(ids, 10, 64)
				if err != nil {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
					return
				}
				if enable != "" {
					Config.Lock.Lock()
					defer Config.Lock.Unlock()
					Enable := ctx.FormValue("enable") == "true"
					if id == -1 {
						Config.CoreConfig.DefaultGroupConfig.Enable = Enable
						_, _ = ctx.JSON(WebResult{
							Code: 1,
							Info: "默认配置保存成功!",
							Data: Config.CoreConfig.GroupConfig[id].Enable,
						})
						err := Config.Save()
						if err != nil {
							log.Println(err)
						}
						return
					}
					if v, ok := Config.CoreConfig.GroupConfig[id]; ok {
						v.Enable = Enable
						Config.CoreConfig.GroupConfig[id] = v
					} else {
						v = Config.CoreConfig.DefaultGroupConfig
						v.Enable = Enable
						Config.CoreConfig.GroupConfig[id] = v
					}
					_, _ = ctx.JSON(WebResult{
						Code: 1,
						Info: "保存成功!",
						Data: Config.CoreConfig.GroupConfig[id].Enable,
					})
					err := Config.Save()
					if err != nil {
						log.Println(err)
					}
					return
				}
				menuData := ctx.FormValue("data[Menu]")
				menuKeyWordData := ctx.FormValue("data[MenuKeyWord]")
				Enable := ctx.FormValue("data[Enable]") == "true"
				ShutUpWord := ctx.FormValue("data[ShutUpWord]")
				Welcome := ctx.FormValue("data[Welcome]")
				AdminUin, _ := strconv.ParseInt(ctx.FormValue("data[AdminUin]"), 10, 64)
				JoinVerifyTime, _ := strconv.Atoi(ctx.FormValue("data[JoinVerifyTime]"))
				JoinAutoShutUpTime, _ := strconv.Atoi(ctx.FormValue("data[JoinAutoShutUpTime]"))
				ShutUpTime, _ := strconv.Atoi(ctx.FormValue("data[ShutUpTime]"))
				JoinVerifyType, _ := strconv.Atoi(ctx.FormValue("data[JoinVerifyType]"))
				Zan := ctx.FormValue("data[Zan]") == "true"
				Bili := ctx.FormValue("data[Bili]") == "true"
				SignIn := ctx.FormValue("data[SignIn]") == "true"
				Job := map[string]Config.Job{}
				for k, v := range ctx.FormValues() {
					//log.Println(k,strings.HasPrefix(k,"data[Job]["),strings.Split(strings.TrimPrefix(k,"data[Job]["),"]"))
					if strings.HasPrefix(k, "data[Job][") {
						if v1 := strings.Split(strings.TrimPrefix(k, "data[Job]["), "]"); len(v1) >= 2 && len(v) >= 1 {
							switch v1[1] {
							case "[Cron":
								v2, _ := Job[v1[0]]
								v2.Cron = v[0]
								Job[v1[0]] = v2
							case "[JobType":
								v2, _ := Job[v1[0]]
								v2.Type, _ = strconv.Atoi(v[0])
								Job[v1[0]] = v2
							case "[Content":
								v2, _ := Job[v1[0]]
								v2.Content = v[0]
								Job[v1[0]] = v2
							}

						}
					}
				}
				Config.Lock.Lock()
				defer Config.Lock.Unlock()

				if id == -1 {
					Config.CoreConfig.DefaultGroupConfig = Config.GroupConfig{BiliUps: map[int64]Config.Up{}, Bili: Bili, Job: Job, JoinVerifyType: JoinVerifyType, Welcome: Welcome, SignIn: SignIn, Zan: Zan, JoinVerifyTime: JoinVerifyTime, JoinAutoShutUpTime: JoinAutoShutUpTime, AdminUin: AdminUin, Menu: menuData, MenuKeyWord: menuKeyWordData, Enable: Enable, ShutUpWord: ShutUpWord, ShutUpTime: ShutUpTime}
					Config.Save()
					_, _ = ctx.JSON(WebResult{
						Code: 1,
						Info: "默认配置，保存成功!",
						Data: nil,
					})
					return
				}
				Config.CoreConfig.GroupConfig[id] = Config.GroupConfig{BiliUps: Config.CoreConfig.GroupConfig[id].BiliUps, Bili: Bili, Job: Job, JoinVerifyType: JoinVerifyType, Welcome: Welcome, SignIn: SignIn, Zan: Zan, JoinVerifyTime: JoinVerifyTime, JoinAutoShutUpTime: JoinAutoShutUpTime, AdminUin: AdminUin, Menu: menuData, MenuKeyWord: menuKeyWordData, Enable: Enable, ShutUpWord: ShutUpWord, ShutUpTime: ShutUpTime}
				Config.Save()
				_, _ = ctx.JSON(WebResult{
					Code: 1,
					Info: "保存成功!",
					Data: nil,
				})
				return
			})
			needAuth.Post("/groupStatus", func(ctx iris.Context) {
				ids := ctx.FormValue("id")
				id, err := strconv.ParseInt(ids, 10, 64)
				if err != nil {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
					return
				}
				Config.Lock.RLock()
				defer Config.Lock.RUnlock()
				if id == -1 {
					_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: Config.CoreConfig.DefaultGroupConfig})
					return
				}
				if v, ok := Config.CoreConfig.GroupConfig[id]; ok {
					_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: v})
					return
				} else {
					_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: Config.CoreConfig.DefaultGroupConfig})
					return
				}
			})
			needAuth.Get("/groups", func(ctx iris.Context) {
				g, err := b.GetGroupList("")
				if err != nil {
					_, _ = ctx.JSON(WebResult{
						Code: 0,
						Info: err.Error(),
						Data: nil,
					})
					return
				}
				_, _ = ctx.JSON(WebResult{
					Code: 1,
					Info: "success",
					Data: g,
				})
			})
			needAuth.Post("/shutUp", func(ctx iris.Context) {
				ids := ctx.FormValue("id")
				uins := ctx.FormValue("uin")
				times := ctx.FormValue("time")
				id, err := strconv.ParseInt(ids, 10, 64)
				if err != nil {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
					return
				}
				uin, err := strconv.ParseInt(uins, 10, 64)
				if err != nil {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
					return
				}
				time1, err := strconv.Atoi(times)
				if err != nil {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
					return
				}
				err = b.SetForbidden(1, time1, id, uin)
				if err != nil {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
					return
				}
				_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: nil})
				return
			})
			needAuth.Post("/kick", func(ctx iris.Context) {
				ids := ctx.FormValue("id")
				uins := ctx.FormValue("uin")
				id, err := strconv.ParseInt(ids, 10, 64)
				if err != nil {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
					return
				}
				uin, err := strconv.ParseInt(uins, 10, 64)
				if err != nil {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
					return
				}
				err = b.KickGroupMember(id, uin)
				if err != nil {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
					return
				}
				_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: nil})
				return
			})
			needAuth.Get("/logout", func(ctx iris.Context) {
				s := sess.Start(ctx)
				s.Set("auth", false)
				s.Clear()
				_, _ = ctx.JSON(WebResult{
					Code: 1,
					Info: "Success",
					Data: nil,
				})
			})
		}
		// App 相关接口实现
		needAppAuth := App.Party("/api/app", requireAppLogin)
		{
			needAppAuth.Get("/status", func(ctx iris.Context) {
				userInfo, _ := b.GetUserInfo(b.QQ)
				ctx.JSON(userInfo)
			})
			needAppAuth.Get("/plugins", func(ctx iris.Context) {
				plugins := []Core.ModuleInfo{}
				log.Info(Core.Modules)
				for _, v := range Core.Modules {
					plugins = append(plugins, v.ModuleInfo())
				}
				ctx.JSON(WebResult{
					Code: 0,
					Info: "success",
					Data: plugins,
				})

			})
		}
		go func() {
			_ = <-Start
			App.Logger().Prefix = "[Web]"
			err := App.Run(iris.Addr(Config.CoreConfig.OPQWebConfig.Host+":"+strconv.Itoa(Config.CoreConfig.OPQWebConfig.Port)), iris.WithoutStartupLog)
			if err != nil {
				log.Println(err)
				return
			}
		}()

	}
	return nil
}

func init() {
	Core.RegisterModule(&Module{})
}

func WhiteGroupList(botQQ int64, packet *OPQBot.GroupMsgPack) {
	if len(Config.CoreConfig.WhiteGroupList) == 0 {
		packet.Next(botQQ, packet)
		return
	}
	isWhite := false
	for _, v := range Config.CoreConfig.WhiteGroupList {
		if v == packet.FromGroupID {
			isWhite = true
			break
		}
	}
	if isWhite {
		packet.Next(botQQ, &packet)
	}
}
func BlackGroupList(botQQ int64, packet *OPQBot.GroupMsgPack) {
	if len(Config.CoreConfig.BlackGroupList) == 0 {
		packet.Next(botQQ, packet)
		return
	}
	isBlack := false
	for _, v := range Config.CoreConfig.WhiteGroupList {
		if v == packet.FromGroupID {
			isBlack = true
			break
		}
	}
	if !isBlack {
		packet.Next(botQQ, packet)
	}
}
func requireAppLogin(ctx iris.Context) {
	ctx.Next()
	return
	s := sess.Start(ctx)
	if s.GetBooleanDefault("appAuth", false) {
		ctx.Next()
	} else {
		_, _ = ctx.JSON(WebResult{Code: 10010, Info: "App 未登录!", Data: nil})
		return
	}
}
func requireAuth(ctx iris.Context) {
	s := sess.Start(ctx)
	if s.GetBooleanDefault("auth", false) {
		ctx.Next()
	} else {
		_, _ = ctx.JSON(WebResult{Code: 10010, Info: "未登录!", Data: nil})
		return
	}
}
func pathIsFile(path string) (isFile bool) {
	isFile = false
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			isFile = true
			break
		}
	}
	return
}
func beforeCsrf(ctx iris.Context) {
	s := sess.Start(ctx)
	//log.Println(s.Get("OPQWebCSRF"))
	if ctx.Method() == "POST" {
		if key := s.GetStringDefault("OPQWebCSRF", ""); key != "" && (ctx.GetHeader("csrfToken") == key || ctx.FormValue("csrfToken") == key) {
			ctx.Next()
		} else {
			// log.Println(key, "-", ctx.FormValue("csrfToken"))
			for _, v := range CsrfWhiteList {
				if strings.HasPrefix(ctx.Path(), v) {
					ctx.Next()
					return
				}
			}

			ctx.StatusCode(419)
			_, _ = ctx.JSON(WebResult{Code: 419, Info: "CSRF Error!", Data: nil})
			return
		}
	} else {
		ctx.Next()
	}
}

var CsrfWhiteList = []string{
	"/github/webhook",
	"/api/kiss",
}
