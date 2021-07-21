package Core

import (
	"OPQBot-QQGroupManager/utils"
	"errors"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/mcoo/OPQBot"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

var Modules = make(map[string]Module)

func init() {
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
func GetLog() *logrus.Logger {
	return log
}
func InitModule(b *Bot) {
	for k, v := range Modules {
		l := log.WithField("Module", k)
		l.Infof("Author: %s - %s", v.ModuleInfo().Author, v.ModuleInfo().Description)
		l.Info("正在载入中")
		err := v.ModuleInit(b, l)
		if err != nil {
			l.Error("导入模块时出错！", err)
		}
		l.Infof("载入成功")
	}
}
func RegisterModule(module Module) error {
	if _, ok := Modules[module.ModuleInfo().Name]; ok {
		return errors.New(module.ModuleInfo().Name + "模块名字已经被注册了")
	} else {
		Modules[module.ModuleInfo().Name] = module
	}
	return nil
}

// Bot 内置了"周期任务","数据库"
type Bot struct {
	OPQBot.BotManager
	BotCronManager utils.BotCron
	Modules        map[string]*Module
}
type ModuleInfo struct {
	Name        string
	Author      string
	Description string
	Version     int
}
type Module interface {
	ModuleInit(bot *Bot, log *logrus.Entry) error
	ModuleInfo() ModuleInfo
}
