package taobao

import (
	"OPQBot-QQGroupManager/Core"
	"encoding/json"
	"fmt"
	"github.com/mcoo/OPQBot"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type Module struct {
	MsgChannel chan OPQBot.GroupMsgPack
	ImgServer  string
}

var (
	log *logrus.Entry
)

type setu struct {
	Title string `json:"title"`
	Pic   string `json:"pic"`
}

func (m *Module) ModuleInfo() Core.ModuleInfo {
	return Core.ModuleInfo{
		Name:        "淘宝买家秀",
		Author:      "bypanghu",
		Description: "生成淘宝买家秀",
		Version:     0,
	}
}

func (m *Module) ModuleInit(b *Core.Bot, l *logrus.Entry) error {
	log = l

	_, err := b.AddEvent(OPQBot.EventNameOnGroupMessage, func(botQQ int64, packet *OPQBot.GroupMsgPack) {
		if packet.FromUserID == botQQ {
			return
		}
		if packet.Content == "买家秀" {
			client := &http.Client{}
			baseUrl := "https://api.vvhan.com/api/tao?type=json"
			req, err := http.NewRequest("GET", baseUrl, nil)
			req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
			if err != nil {
				panic(err)
			}
			response, _ := client.Do(req)
			defer response.Body.Close()
			s, err := ioutil.ReadAll(response.Body)
			var res setu
			json.Unmarshal(s, &res)
			b.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeGroup,
				ToUserUid:  packet.FromGroupID,
				Content: OPQBot.SendTypePicMsgByUrlContent{
					Content: fmt.Sprintf("%s\n图片地址为：%s\n30s自动撤回\n%s", res.Title, res.Pic, OPQBot.MacroId()),
					PicUrl:  res.Pic,
					Flash:   false,
				},
				CallbackFunc: func(Code int, Info string, record OPQBot.MyRecord) {
					time.Sleep(30 * time.Second)
					err := b.ReCallMsg(packet.FromGroupID, record.MsgRandom, record.MsgSeq)
					if err != nil {
						log.Warn(err)
					}
				},
			})
		}

	})
	return err
}

func init() {
	Core.RegisterModule(&Module{})
}
