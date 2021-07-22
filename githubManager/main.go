package githubManager

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/Core"
	"OPQBot-QQGroupManager/GroupManager"
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/mcoo/OPQBot"
	"strings"
)

type Module struct {
}

var log *logrus.Entry

func (m *Module) ModuleInfo() Core.ModuleInfo {
	return Core.ModuleInfo{
		Name:          "Github订阅姬",
		Author:        "enjoy",
		Description:   "",
		Version:       0,
		RequireModule: []string{"群管理插件"},
	}
}
func (m *Module) ModuleInit(b *Core.Bot, l *logrus.Entry) error {
	log = l
	g := NewManager(GroupManager.App, b)
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
		if packet.Content == "本群Github" {
			githubs := "本群订阅Github仓库\n"
			list := g.GetGroupSubList(packet.FromGroupID)
			if len(list) == 0 {
				b.SendGroupTextMsg(packet.FromGroupID, "本群没有订阅Github仓库")
				return
			}
			for k, _ := range list {
				githubs += fmt.Sprintf("%s \n", k)
			}
			b.SendGroupTextMsg(packet.FromGroupID, githubs)
			return
		}
		if len(cm) == 2 && cm[0] == "取消订阅Github" {
			err := g.DelRepo(cm[1], packet.FromGroupID)
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			b.SendGroupTextMsg(packet.FromGroupID, "取消订阅成功!")
			return
		}
		if len(cm) == 2 && cm[0] == "订阅Github" {
			b.SendGroupTextMsg(packet.FromGroupID, "请私聊我发送该仓库的Webhook Secret!")
			err := s.Set("github", cm[1])
			if err != nil {
				log.Println(err)
				return
			}
			err = s.Set("github_groupId", packet.FromGroupID)
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	})
	if err != nil {
		return err
	}

	err = b.AddEvent(OPQBot.EventNameOnFriendMessage, func(qq int64, packet *OPQBot.FriendMsgPack) {
		s := b.Session.SessionStart(packet.FromUin)
		if v, err := s.GetString("github"); err == nil {
			groupidI, err := s.Get("github_groupId")
			if err != nil {
				b.SendFriendTextMsg(packet.FromUin, err.Error())
				return
			}
			groupId, ok := groupidI.(int64)
			if !ok {
				b.SendFriendTextMsg(packet.FromUin, "内部错误")
				return
			}
			err = g.AddRepo(v, packet.Content, groupId)
			if err != nil {
				b.SendFriendTextMsg(packet.FromUin, err.Error())
				s.Delete("github")
				s.Delete("github_groupId")
				return
			}
			b.SendFriendTextMsg(packet.FromUin, "成功!")
			s.Delete("github")
			s.Delete("github_groupId")
		}
	})
	if err != nil {
		log.Println(err)
	}
	return nil
}

func init() {
	Core.RegisterModule(&Module{})
}
