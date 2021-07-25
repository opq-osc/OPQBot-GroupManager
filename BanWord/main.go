package BanWord

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/Core"
	"github.com/mcoo/OPQBot"
	"log"
	"regexp"
)

func Hook(b *Core.Bot) error {
	b.AddEvent(OPQBot.EventNameOnGroupMessage, func(botQQ int64, packet *OPQBot.GroupMsgPack) {
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
		banQQ := Config.CoreConfig.BanQQ
		Config.Lock.RUnlock()

		for _, v := range banQQ {
			if packet.FromUserID == v {
				packet.Ban = true
				return
			}
		}

		if !c.Enable {
			return
		}
		if m, err := regexp.MatchString(c.ShutUpWord, packet.Content); err != nil {
			log.Println(err)
			return
		} else if m {
			err := packet.Bot.ReCallMsg(packet.FromGroupID, packet.MsgRandom, packet.MsgSeq)
			if err != nil {
				log.Println(err)
			}
			err = packet.Bot.SetForbidden(1, c.ShutUpTime, packet.FromGroupID, packet.FromUserID)
			if err != nil {
				log.Println(err)
			}
			packet.Bot.SendGroupTextMsg(packet.FromGroupID, OPQBot.MacroAt([]int64{packet.FromUserID})+"触发违禁词")
			packet.Ban = true
			return
		}
	})
	return nil
}
