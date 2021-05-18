package main

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/methods"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/mcoo/OPQBot"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	version = "Ver.0.0.1"
	sess    *sessions.Sessions
)

type WebResult struct {
	Code int         `json:"code"`
	Info string      `json:"info"`
	Data interface{} `json:"data"`
}

func main() {
	log.Println("QQ Group Manager✈️" + version)
	b := OPQBot.NewBotManager(Config.CoreConfig.OPQBotConfig.QQ, Config.CoreConfig.OPQBotConfig.Url)
	err := b.AddEvent(OPQBot.EventNameOnDisconnected, func() {
		log.Println("断开服务器")
	})
	if err != nil {
		log.Println(err)
	}
	// 黑名单优先级高于白名单
	err = b.AddEvent(OPQBot.EventNameOnGroupMessage, BlackGroupList, WhiteGroupList, func(botQQ int64, packet *OPQBot.GroupMsgPack) {
		if packet.FromUserID == botQQ {
			return
		}
		Config.Lock.RLock()
		defer Config.Lock.RUnlock()
		var c Config.GroupConfig
		if v, ok := Config.CoreConfig.GroupConfig[packet.FromGroupID]; ok {
			c = v
		} else {
			c = Config.CoreConfig.DefaultGroupConfig
		}
		if !c.Enable {
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
		if m, err := regexp.MatchString(c.ShutUpWord, packet.Content); err != nil {
			log.Println(err)
			return
		} else if m {
			err := b.ReCallMsg(packet.FromGroupID, packet.MsgRandom, packet.MsgSeq)
			if err != nil {
				log.Println(err)
			}
			err = b.SetForbidden(1, c.ShutUpTime, packet.FromGroupID, packet.FromUserID)
			if err != nil {
				log.Println(err)
			}
			return
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
			if v,ok := Config.CoreConfig.UserData[packet.FromUserID];ok {
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
			}else{
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
			if v,ok := Config.CoreConfig.UserData[packet.FromUserID];ok {
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
			}else{
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
		}
		if packet.Content == "积分" {
			if v,ok := Config.CoreConfig.UserData[packet.FromUserID];ok {
				b.Send(OPQBot.SendMsgPack{
					SendToType:   OPQBot.SendToTypeGroup,
					ToUserUid:    packet.FromGroupID,
					Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "你的积分为"+strconv.Itoa(v.Count)},
					CallbackFunc: nil,
				})
			}else{
				b.Send(OPQBot.SendMsgPack{
					SendToType:   OPQBot.SendToTypeGroup,
					ToUserUid:    packet.FromGroupID,
					Content:      OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "你的积分为0"},
					CallbackFunc: nil,
				})
			}
		}
 	})
	if err != nil {
		log.Println(err)
	}
	err = b.AddEvent(OPQBot.EventNameOnGroupJoin, func(botQQ int64, packet *OPQBot.GroupJoinPack) {

	})
	err = b.AddEvent(OPQBot.EventNameOnConnected, func() {
		log.Println("连接服务器成功")
	})
	if err != nil {
		log.Println(err)
	}
	err = b.Start()
	if err != nil {
		log.Println(err)
	}
	if Config.CoreConfig.OPQWebConfig.Enable {
		log.Println("启动Web 😊")
		go func() {
			app := iris.New()
			Config.Lock.Lock()
			sess = sessions.New(sessions.Config{Cookie: "OPQWebSession"})
			if Config.CoreConfig.OPQWebConfig.CSRF == "" {
				Config.CoreConfig.OPQWebConfig.CSRF = RandomString(32)
				err := Config.Save()
				if err != nil {
					log.Println(err)
				}
			}

			Config.Lock.Unlock()
			app.Use(beforeCsrf)
			app.Use(sess.Handler())
			app.WrapRouter(func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
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
					if r.URL.Path[0:4] != "/api" {
						if !pathIsFile(path) {
							r.URL.Path = "/"
						}
					}
				}

				router.ServeHTTP(w, r)
			})
			app.HandleDir("/", iris.Dir("./Web/dist/spa"))
			app.Get("/api/status", func(ctx iris.Context) {
				s := sess.Start(ctx)
				salt := int(time.Now().Unix())
				keyTmp := methods.Md5V(strconv.Itoa(salt + rand.Intn(100)))
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
			app.Post("/api/login", func(ctx iris.Context) {
				username := ctx.FormValue("username")
				password := ctx.FormValue("password")
				Config.Lock.RLock()
				defer Config.Lock.RUnlock()
				if username == Config.CoreConfig.OPQWebConfig.Username && password == methods.Md5V(Config.CoreConfig.OPQWebConfig.Password) {
					s := sess.Start(ctx)
					s.Set("auth", true)
					_, _ = ctx.JSON(WebResult{Code: 1, Info: "登录成功", Data: nil})
					return
				} else {
					_, _ = ctx.JSON(WebResult{Code: 0, Info: "用户名密码错误!", Data: nil})
					return
				}

			})
			needAuth := app.Party("/api/admin", requireAuth)
			{
				needAuth.Post("/getGroupMember", func(ctx iris.Context) {
					ids := ctx.FormValue("id")
					id, err := strconv.ParseInt(ids, 10, 64)
					if id == -1 {
						_, _ = ctx.JSON(WebResult{Code: 1, Info: "success", Data: []int{}})
						return
					}
					if err != nil {
						_, _ = ctx.JSON(WebResult{Code: 0, Info: err.Error(), Data: nil})
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
					SignIn := ctx.FormValue("data[SignIn]") == "true"
					Job := map[string]Config.Job{}
					for k,v := range ctx.FormValues() {
						//log.Println(k,strings.HasPrefix(k,"data[Job]["),strings.Split(strings.TrimPrefix(k,"data[Job]["),"]"))
						if strings.HasPrefix(k,"data[Job][") {
							if v1 := strings.Split(strings.TrimPrefix(k,"data[Job]["),"]");len(v1) >=2 && len(v)>=1 {
								switch v1[1] {
								case "[Cron":
									v2,_ := Job[v1[0]]
									v2.Cron = v[0]
									Job[v1[0]] = v2
								case "[JobType":
									v2,_ := Job[v1[0]]
									v2.Type,_ = strconv.Atoi(v[0])
									Job[v1[0]] = v2
								case "[Content":
									v2,_ := Job[v1[0]]
									v2.Content=v[0]
									Job[v1[0]] = v2
								}
								
							}
						}
					}
					Config.Lock.Lock()
					defer Config.Lock.Unlock()

					if id == -1 {
						Config.CoreConfig.DefaultGroupConfig = Config.GroupConfig{Job: Job, JoinVerifyType: JoinVerifyType, Welcome: Welcome, SignIn: SignIn, Zan: Zan, JoinVerifyTime: JoinVerifyTime, JoinAutoShutUpTime: JoinAutoShutUpTime, AdminUin: AdminUin, Menu: menuData, MenuKeyWord: menuKeyWordData, Enable: Enable, ShutUpWord: ShutUpWord, ShutUpTime: ShutUpTime}
						Config.Save()
						_, _ = ctx.JSON(WebResult{
							Code: 1,
							Info: "默认配置，保存成功!",
							Data: nil,
						})
						return
					}
					Config.CoreConfig.GroupConfig[id] = Config.GroupConfig{Job: Job,JoinVerifyType: JoinVerifyType, Welcome: Welcome, SignIn: SignIn, Zan: Zan, JoinVerifyTime: JoinVerifyTime, JoinAutoShutUpTime: JoinAutoShutUpTime, AdminUin: AdminUin, Menu: menuData, MenuKeyWord: menuKeyWordData, Enable: Enable, ShutUpWord: ShutUpWord, ShutUpTime: ShutUpTime}
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
			app.Logger().Prefix = "[Web]"
			err := app.Run(iris.Addr(Config.CoreConfig.OPQWebConfig.Host+":"+strconv.Itoa(Config.CoreConfig.OPQWebConfig.Port)), iris.WithoutStartupLog)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
	b.Wait()
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
func requireAuth(ctx iris.Context) {
	s := sess.Start(ctx)
	if s.GetBooleanDefault("auth", false) {
		ctx.Next()
	} else {
		_, _ = ctx.JSON(WebResult{Code: 1, Info: "未登录!", Data: nil})
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
			log.Println(key, "-", ctx.FormValue("csrfToken"))
			ctx.StatusCode(419)
			_, _ = ctx.Text("CSRF Error!")
			return
		}
	} else {
		ctx.Next()
	}
}
