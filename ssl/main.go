package ssl

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/Core"
	"OPQBot-QQGroupManager/utils"
	"fmt"
	"github.com/mcoo/OPQBot"
	"github.com/sirupsen/logrus"
	"strings"
)

type Module struct {
}

var log *logrus.Entry

func (m *Module) ModuleInfo() Core.ModuleInfo {
	return Core.ModuleInfo{
		Name:        "SSL扫描插件",
		Author:      "enjoy",
		Description: "",
		Version:     0,
	}
}
func (m *Module) ModuleInit(b *Core.Bot, l *logrus.Entry) error {
	log = l
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
		if len(cm) == 2 && cm[0] == "SSL" {
			ssl, err := utils.SSLStatus("github.com")
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			content := "SSL证书建议:"
			for _, v := range ssl.Data.Status.Suggests {
				content += fmt.Sprintf("\n- %s", v.Tip)
			}
			content += "\n具体信息:"
			for _, v := range ssl.Data.Status.Protocols {
				content += fmt.Sprintf("\n- %s %t", v.Name, v.Support)
			}
			for _, v := range ssl.Data.Status.ProtocolDetail {
				content += fmt.Sprintf("\n- %s %t", v.Name, v.Support)
			}
			b.SendGroupTextMsg(packet.FromGroupID, OPQBot.MacroAt([]int64{packet.FromUserID})+content)
			return
		}
		if len(cm) == 2 && cm[0] == "DNS" {
			dns, err := utils.DnsQuery(cm[1])
			if err != nil {
				b.SendGroupTextMsg(packet.FromGroupID, err.Error())
				return
			}
			content := fmt.Sprintf("[%s]\n", cm[1])
			if dns.Data.Num86[0].Answer.Error == "" {
				content += fmt.Sprintf("中国: %s [%s] %ss\n", dns.Data.Num86[0].Answer.Records[0].Value, dns.Data.Num86[0].Answer.Records[0].IPLocation, dns.Data.Num86[0].Answer.TimeConsume)
			} else {
				content += "中国: " + dns.Data.Num86[0].Answer.Error + "\n"
			}
			if dns.Data.Num01[0].Answer.Error == "" {
				content += fmt.Sprintf("美国: %s [%s] %ss\n", dns.Data.Num01[0].Answer.Records[0].Value, dns.Data.Num01[0].Answer.Records[0].IPLocation, dns.Data.Num01[0].Answer.TimeConsume)
			} else {
				content += "美国: " + dns.Data.Num01[0].Answer.Error + "\n"
			}
			if dns.Data.Num852[0].Answer.Error == "" {
				content += fmt.Sprintf("香港: %s [%s] %ss", dns.Data.Num852[0].Answer.Records[0].Value, dns.Data.Num852[0].Answer.Records[0].IPLocation, dns.Data.Num852[0].Answer.TimeConsume)
			} else {
				content += "香港: " + dns.Data.Num852[0].Answer.Error
			}
			b.SendGroupTextMsg(packet.FromGroupID, OPQBot.MacroAt([]int64{packet.FromUserID})+content)
			return
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func init() {
	Core.RegisterModule(&Module{})
}
