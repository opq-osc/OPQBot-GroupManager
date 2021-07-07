package main

import (
	"OPQBot-QQGroupManager/BanWord"
	_ "OPQBot-QQGroupManager/Bili"
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/Core"
	_ "OPQBot-QQGroupManager/GroupManager"
	"OPQBot-QQGroupManager/androidDns"
	_ "OPQBot-QQGroupManager/genAndYiqin"
	_ "OPQBot-QQGroupManager/githubManager"
	//_ "OPQBot-QQGroupManager/steam"
	"OPQBot-QQGroupManager/utils"

	_ "github.com/go-playground/webhooks/v6/github"
	"github.com/mcoo/OPQBot"

	"log"
)

var (
	version = "unknown"
	date    = "none"
)

func main() {
	log.Println("QQ Group Manager -️" + version + " 编译时间 " + date)
	androidDns.SetDns()
	go CheckUpdate()
	b := Core.Bot{BotManager: OPQBot.NewBotManager(Config.CoreConfig.OPQBotConfig.QQ, Config.CoreConfig.OPQBotConfig.Url)}
	err := b.AddEvent(OPQBot.EventNameOnDisconnected, func() {
		log.Println("断开服务器")
	})
	if err != nil {
		log.Println(err)
	}
	b.BotCronManager = utils.NewBotCronManager()
	b.BotCronManager.Start()
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
	_ = BanWord.Hook(&b)
	for _, v := range Core.Modules {
		err := v.ModuleInit(&b)
		if err != nil {
			log.Println("导入模块时出错！", err)
		}
	}
	b.Wait()
}
