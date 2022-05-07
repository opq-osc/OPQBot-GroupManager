package Core

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/utils"
	"errors"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/mcoo/OPQBot"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"runtime"
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

var hasLoad []string

func hasLoaded(infoName string) bool {
	tmp := false
	for _, v1 := range hasLoad {
		if v1 == infoName {
			tmp = true
			break
		}
	}
	return tmp
}
func InitModule(b *Bot) {
	for _, v := range Modules {
		if err := startModule(b, v); err != nil {
			log.Error(err)
		}
	}
}
func startModule(b *Bot, module Module) error {
	info := module.ModuleInfo()
	if hasLoaded(info.Name) {
		return nil
	}
	l := log.WithField("Module", info.Name)
	for _, v2 := range info.RequireModule {
		if !hasLoaded(v2) {
			if v, ok := Modules[v2]; ok {
				err := startModule(b, v)
				if err != nil {
					return err
				}
			} else {
				return errors.New(fmt.Sprintf("缺少依赖%s 导入失败\n", v2))
			}
		}
	}

	l.Infof("Author: %s - %s", info.Author, info.Description)
	l.Info("正在载入中")
	err := module.ModuleInit(b, l)
	if err != nil {
		l.Error("导入模块时出错！", err)
	}
	l.Infof("载入成功")
	hasLoad = append(hasLoad, info.Name)
	return nil
}
func RegisterModule(module Module) error {
	Config.Lock.RLock()
	disable := Config.CoreConfig.DisableModule
	Config.Lock.RUnlock()
	for _, v := range disable {
		if v == module.ModuleInfo().Name {
			log.Warn(module.ModuleInfo().Name + "模块已被禁止载入")
			return errors.New("模块已被禁止载入")
		}
	}
	if _, ok := Modules[module.ModuleInfo().Name]; ok {
		log.Error(module.ModuleInfo().Name + "模块名字已经被注册了")
		return errors.New(module.ModuleInfo().Name + "模块名字已经被注册了")
	} else {
		Modules[module.ModuleInfo().Name] = module
	}
	return nil
}

var lastTotalFreed uint64

func (b *Bot) PrintMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Alloc = %v TotalAlloc = %v  Just Freed = %v Sys = %v NumGC = %v\n",
		m.Alloc/1024, m.TotalAlloc/1024, ((m.TotalAlloc-m.Alloc)-lastTotalFreed)/1024, m.Sys/1024, m.NumGC)
	lastTotalFreed = m.TotalAlloc - m.Alloc
}

// Bot 内置了"周期任务","数据库"
type Bot struct {
	*OPQBot.BotManager
	BotCronManager utils.BotCron
	Modules        map[string]*Module
	DB             *gorm.DB
}
type ModuleInfo struct {
	Name          string
	Author        string
	Description   string
	Version       int
	RequireModule []string
}
type Module interface {
	ModuleInit(bot *Bot, log *logrus.Entry) error
	ModuleInfo() ModuleInfo
}
