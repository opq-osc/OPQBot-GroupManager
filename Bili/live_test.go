package bili

import (
	"github.com/gaoyanpao/biliLiveHelper"
	"log"
	"os"
	"testing"
)

func TestLive(t *testing.T) {
	client := biliLiveHelper.NewClient(650)
	if client == nil {
		os.Exit(1)
	}
	client.PrintRoomInfo()
	client.RegHandleFunc(biliLiveHelper.CmdDanmuMsg, func(ctx *biliLiveHelper.Context) bool {
		data := ctx.Msg
		log.Printf("[弹幕]<%v>%s", data.Get("info").GetIndex(2).GetIndex(1).MustString(), data.Get("info").GetIndex(1).MustString())
		return false
	})
	client.RegHandleFunc(biliLiveHelper.CmdWelcome, func(ctx *biliLiveHelper.Context) bool {
		data := ctx.Msg
		log.Printf("[欢迎]%v", data.Get("data"))
		return false
	})
	client.RegHandleFunc(biliLiveHelper.CmdSendGift, func(ctx *biliLiveHelper.Context) bool {
		data := ctx.Msg
		log.Printf("[礼物]%v", data.Get("data"))
		return false
	})
	client.StartListen()
}
