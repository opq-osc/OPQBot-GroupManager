package kiss

import (
	"OPQBot-QQGroupManager/Core"
	mydraw "OPQBot-QQGroupManager/draw"
	"bytes"
	"embed"
	"encoding/json"
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/gg"
	"github.com/sirupsen/logrus"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"strconv"
	"strings"
)

type Module struct {
}

var (
	log *logrus.Entry
)

//go:embed static
var static embed.FS
var (
	OPERATOR_X = []int{92, 135, 84, 80, 155, 60, 50, 98, 35, 38, 70, 84, 75}
	OPERATOR_Y = []int{64, 40, 105, 110, 82, 96, 80, 55, 65, 100, 80, 65, 65}
	TARGET_X   = []int{58, 62, 42, 50, 56, 18, 28, 54, 46, 60, 35, 20, 40}
	TARGET_Y   = []int{90, 95, 100, 100, 100, 120, 110, 100, 100, 100, 115, 120, 96}
)

func (m *Module) ModuleInfo() Core.ModuleInfo {
	return Core.ModuleInfo{
		Name:          "Kiss",
		Author:        "enjoy",
		Description:   "嘿嘿嘿",
		Version:       0,
		RequireModule: []string{"群管理插件"},
	}
}

type AtMsg struct {
	Content string `json:"Content"`
	UserExt []struct {
		QQNick string `json:"QQNick"`
		QQUID  int64  `json:"QQUid"`
	} `json:"UserExt"`
	UserID []int64 `json:"UserID"`
}

var pics []image.Image

func isInPalette(p color.Palette, c color.Color) int {
	ret := -1
	for i, v := range p {
		if v == c {
			return i
		}
	}
	return ret
}
func getPalette(m image.Image) color.Palette {
	p := color.Palette{color.RGBA{}}
	p9 := color.Palette(palette.WebSafe)
	b := m.Bounds()
	black := false
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := m.At(x, y)
			cc := p9.Convert(c)
			if cc == p9[0] {
				black = true
			}
			if isInPalette(p, cc) == -1 {
				p = append(p, cc)
			}
		}
	}
	if len(p) < 256 && black == true {
		p[0] = color.RGBA{} // transparent
		p = append(p, p9[0])
	}
	return p
}
func DrawKiss(ava1, ava2 image.Image) []byte {
	var disposals []byte
	var images []*image.Paletted
	var delays []int

	for i, v := range pics {
		bg := gg.NewContextForImage(v)
		bg.DrawImage(ava2, TARGET_X[i], TARGET_Y[i])
		bg.DrawImage(ava1, OPERATOR_X[i], OPERATOR_Y[i])
		img := bg.Image()
		cp := getPalette(img)
		disposals = append(disposals, gif.DisposalNone) //透明图片需要设置
		p := image.NewPaletted(image.Rect(0, 0, bg.Height(), bg.Width()), cp)
		draw.Draw(p, p.Bounds(), img, image.Point{}, draw.Src)
		images = append(images, p)
		delays = append(delays, 0)

	}
	g := &gif.GIF{
		Image:     images,
		Delay:     delays,
		LoopCount: 0,
		Disposal:  disposals,
	}
	buf := new(bytes.Buffer)
	err := gif.EncodeAll(buf, g)
	if err != nil {
		log.Error(err)
	}
	return buf.Bytes()
}
func (m *Module) ModuleInit(b *Core.Bot, l *logrus.Entry) error {
	log = l
	for i := 1; i <= 13; i++ {
		picBytes, err := static.ReadFile("static/" + strconv.Itoa(i) + ".png")
		if err != nil {
			continue
		}
		pic, err := png.Decode(bytes.NewBuffer(picBytes))
		if err != nil {
			continue
		}
		pics = append(pics, pic)
	}
	//GroupManager.App.Post("/api/kiss", func(ctx iris.Context) {
	//	srcPic := ctx.FormValueDefault("src","")
	//	targetPic := ctx.FormValueDefault("target","")
	//	if srcPic == "" || targetPic == "" {
	//		ctx.JSON(GroupManager.WebResult{
	//			Code: 1,
	//			Info: "error",
	//			Data: "字段不匹配",
	//		})
	//		return
	//	}
	//	avaTmp1, err := mydraw.GetAvatar(srcPic)
	//	if err != nil {
	//		ctx.JSON(GroupManager.WebResult{
	//			Code: 1,
	//			Info: "error",
	//			Data: err.Error(),
	//		})
	//		return
	//	}
	//	avaTmp2, err := mydraw.GetAvatar(targetPic)
	//	if err != nil {
	//		ctx.JSON(GroupManager.WebResult{
	//			Code: 1,
	//			Info: "error",
	//			Data: err.Error(),
	//		})
	//		return
	//	}
	//	gifPic := DrawKiss(mydraw.DrawCircle(avaTmp1, 40), mydraw.DrawCircle(avaTmp2, 50))
	//	jsonPic,err	 := ctx.URLParamBool("json")
	//	if err != nil {
	//		jsonPic = true
	//	}
	//	if jsonPic {
	//		ctx.JSON(GroupManager.WebResult{
	//			Code: 0,
	//			Info: "success",
	//			Data: base64.StdEncoding.EncodeToString(gifPic),
	//		})
	//	} else {
	//		ctx.StatusCode(200)
	//		ctx.Header("content-type", "image/gif")
	//		ctx.Write(gifPic)
	//	}
	//
	//})
	_, err := b.AddEvent(OPQBot.EventNameOnGroupMessage, func(qq int64, packet *OPQBot.GroupMsgPack) {
		if packet.FromUserID != b.QQ {
			if packet.MsgType == "AtMsg" && strings.Contains(packet.Content, "kiss") {
				var atInfo AtMsg
				if json.Unmarshal([]byte(packet.Content), &atInfo) == nil {
					if len(atInfo.UserID) == 1 {
						avaTmp1, err := mydraw.GetAvatar("http://q1.qlogo.cn/g?b=qq&s=640&nk=" + strconv.FormatInt(packet.FromUserID, 10))
						if err != nil {
							log.Error(err)
							return
						}
						avaTmp2, err := mydraw.GetAvatar("http://q1.qlogo.cn/g?b=qq&s=640&nk=" + strconv.FormatInt(atInfo.UserID[0], 10))
						if err != nil {
							log.Error(err)
							return
						}
						gifPic := DrawKiss(mydraw.DrawCircle(avaTmp1, 40), mydraw.DrawCircle(avaTmp2, 50))
						b.SendGroupPicMsg(packet.FromGroupID, "", gifPic)
					} else if len(atInfo.UserID) >= 2 {
						avaTmp1, err := mydraw.GetAvatar("http://q1.qlogo.cn/g?b=qq&s=640&nk=" + strconv.FormatInt(atInfo.UserID[0], 10))
						if err != nil {
							log.Error(err)
							return
						}
						avaTmp2, err := mydraw.GetAvatar("http://q1.qlogo.cn/g?b=qq&s=640&nk=" + strconv.FormatInt(atInfo.UserID[1], 10))
						if err != nil {
							log.Error(err)
							return
						}
						gifPic := DrawKiss(mydraw.DrawCircle(avaTmp1, 40), mydraw.DrawCircle(avaTmp2, 50))
						b.SendGroupPicMsg(packet.FromGroupID, "", gifPic)
					}
				}
			}
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
