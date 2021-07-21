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
	"github.com/sirupsen/logrus"

	//_ "OPQBot-QQGroupManager/steam"
	"OPQBot-QQGroupManager/utils"

	_ "github.com/go-playground/webhooks/v6/github"
	"github.com/mcoo/OPQBot"
)

var (
	version = "unknown"
	date    = "none"
	log     *logrus.Logger
)

func main() {
	log = Core.GetLog()
	if Config.CoreConfig.LogLevel != "" {
		switch Config.CoreConfig.LogLevel {
		case "info":
			log.SetLevel(logrus.InfoLevel)
		case "debug":
			log.SetLevel(logrus.DebugLevel)
		case "warn":
			log.SetLevel(logrus.WarnLevel)
		case "error":
			log.SetLevel(logrus.ErrorLevel)
		}

	}
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
	b.DB = Config.DB
	err = b.AddEvent(OPQBot.EventNameOnConnected, func() {
		log.Println("连接服务器成功")
	})
	if err != nil {
		log.Println(err)
	}
	_ = BanWord.Hook(&b)
	Core.InitModule(&b)
	err = b.Start()
	if err != nil {
		log.Error(err)
	}
	b.Wait()
}
