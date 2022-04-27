package main

import (
	"OPQBot-QQGroupManager/BanWord"
	_ "OPQBot-QQGroupManager/Bili"
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/Core"
	"OPQBot-QQGroupManager/GroupManager"
	"OPQBot-QQGroupManager/androidDns"
	_ "OPQBot-QQGroupManager/genAndYiqin"
	_ "OPQBot-QQGroupManager/githubManager"
	_ "OPQBot-QQGroupManager/kiss"
	_ "OPQBot-QQGroupManager/setu"
	_ "OPQBot-QQGroupManager/taobao"
	_ "OPQBot-QQGroupManager/wordCloud"
	"bytes"
	_ "embed"
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"

	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"

	//_ "OPQBot-QQGroupManager/steam"
	"OPQBot-QQGroupManager/utils"

	_ "github.com/go-playground/webhooks/v6/github"
	"github.com/mcoo/OPQBot"
)

//go:embed logo.txt
var logo string

var (
	version = "unknown"
	date    = "none"
	log     *logrus.Logger
)

func init() {
	rand.Seed(time.Now().Unix())
}
func main() {
	isEnabled := true
	isColorEnabled := true
	banner.Init(colorable.NewColorableStdout(), isEnabled, isColorEnabled, bytes.NewBufferString(strings.ReplaceAll(logo, "{{ .version }}", version)))
	log = Core.GetLog()
	if Config.CoreConfig.LogLevel != "" {
		switch Config.CoreConfig.LogLevel {
		case "info":
			log.SetLevel(logrus.InfoLevel)
		case "debug":
			log.SetLevel(logrus.DebugLevel)
			log.SetReportCaller(true)
		case "warn":
			log.SetLevel(logrus.WarnLevel)
		case "error":
			log.SetLevel(logrus.ErrorLevel)
		}

	}
	if Config.CoreConfig.Debug {
		log.Warn("注意当前处于DEBUG模式，会开放25569端口，如果你不清楚请关闭DEBUG，因为这样可能泄漏你的信息！😥")
		go func() {
			ip := ":25569"
			if err := http.ListenAndServe(ip, nil); err != nil {
				fmt.Printf("start pprof failed on %s\n", ip)
			}
		}()
	}
	log.Println("QQ Group Manager -️" + version + " 编译时间 " + date)
	androidDns.SetDns()
	go CheckUpdate()
	b := Core.Bot{Modules: map[string]*Core.Module{}, BotManager: OPQBot.NewBotManager(Config.CoreConfig.OPQBotConfig.QQ, Config.CoreConfig.OPQBotConfig.Url)}
	_, err := b.AddEvent(OPQBot.EventNameOnDisconnected, func() {
		log.Println("断开服务器")
	})
	if err != nil {
		log.Println(err)
	}
	b.BotCronManager = utils.NewBotCronManager()
	b.BotCronManager.Start()
	b.DB = Config.DB
	_, err = b.AddEvent(OPQBot.EventNameOnConnected, func() {
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
	GroupManager.Start <- struct{}{}
	b.Wait()
}
