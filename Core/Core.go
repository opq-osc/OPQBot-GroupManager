package Core

import (
	"OPQBot-QQGroupManager/utils"
	"github.com/mcoo/OPQBot"
)

var Modules []Module

func RegisterModule(module Module) {
	Modules = append(Modules, module)
}

type Bot struct {
	OPQBot.BotManager
	BotCronManager utils.BotCron
}

type Module interface {
	ModuleInit(bot *Bot) error
}
