package draw

import (
	"OPQBot-QQGroupManager/Core"
	"OPQBot-QQGroupManager/GroupManager/QunInfo"
	"bytes"
	"crypto/rand"
	_ "embed"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/mcoo/requests"
	"github.com/sirupsen/logrus"
	"image"
	_ "image/jpeg"
	"image/png"
	"math/big"

	"github.com/mcoo/gg"
)

//go:embed techno-hideo-1.ttf
var techno []byte

//go:embed bg.png
var bg []byte

//go:embed AlibabaPuHuiTi-2-55-Regular.ttf
var AliFont []byte

var log *logrus.Logger

func init() {
	log = Core.GetLog()
}
func GetAvatar(avatar string) (image.Image, error) {
	res, err := requests.Get(avatar)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewBuffer(res.Content()))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func DrawGroupInfo(GroupInfo QunInfo.GroupInfoResult, GroupMemberInfo QunInfo.GroupMembersResult) ([]byte, error) {
	bgimg, err := png.Decode(bytes.NewReader(bg))
	if err != nil {
		return nil, err
	}
	dc := gg.NewContextForImage(bgimg)
	err = dc.LoadFontFaceFromBytes(AliFont, 16)
	if err != nil {
		return nil, err
	}
	dc.SetRGB(200, 200, 200)
	a := 0
	for _, i := range GroupMemberInfo.Data.SpeakRank {
		if a >= 4 {
			break
		}
		ava, err := GetAvatar(i.Avatar)
		if err != nil {
			log.Error()
			continue
		}
		ava = imaging.Resize(ava, 64, 64, imaging.Lanczos)
		c := gg.NewContext(64, 64)
		// 画圆形
		c.DrawCircle(32, 32, 32)
		// 对画布进行裁剪
		c.Clip()
		c.DrawImage(ava, 0, 0)
		dc.DrawImage(c.Image(), 16, 136+101*a)

		text := TruncateText(dc, fmt.Sprintf("%s (%s)", i.Nickname, i.Uin), 252)
		dc.DrawString(text, 96, float64(158+101*a))
		text = TruncateText(dc, fmt.Sprintf("活跃度 %d    发言条数 %d", i.Active, i.MsgCount), 252)
		dc.DrawString(text, 96, float64(196+101*a))
		a += 1
	}
	dc.LoadFontFaceFromBytes(AliFont, 11)
	dc.SetRGB(150, 150, 150)
	log.Println(GroupInfo.Data.ActiveData)
	text := ""
	if len(GroupInfo.Data.ActiveData.DataList) == 0 {
		text = TruncateText(dc, fmt.Sprintf("昨日活跃人数null  消息条数null  加群null人  退群null人  申请入群null人"), 390)
	} else {
		text = TruncateText(dc, fmt.Sprintf("昨日活跃人数%d  消息条数%d  加群%d人  退群%d人  申请入群%d人",
			GroupInfo.Data.ActiveData.DataList[len(GroupInfo.Data.ActiveData.DataList)-1].Number,
			GroupInfo.Data.MsgInfo.DataList[len(GroupInfo.Data.MsgInfo.DataList)-1].Number,
			GroupInfo.Data.JoinData.DataList[len(GroupInfo.Data.JoinData.DataList)-1].Number,
			GroupInfo.Data.ExitData.DataList[len(GroupInfo.Data.ExitData.DataList)-1].Number,
			GroupInfo.Data.ApplyData.DataList[len(GroupInfo.Data.ApplyData.DataList)-1].Number,
		), 390)
	}

	dc.DrawString(text, 6, 526)
	// 226- 141 + 141 -126
	buf := new(bytes.Buffer)
	err = png.Encode(buf, dc.Image())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
func TruncateText(dc *gg.Context, originalText string, maxTextWidth float64) string {
	tmpStr := ""
	result := make([]rune, 0)
	for _, r := range originalText {
		tmpStr = tmpStr + string(r)
		w, _ := dc.MeasureString(tmpStr)
		if w > maxTextWidth {
			if len(tmpStr) <= 1 {
				return ""
			} else {
				break
			}
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
func Draw6Number() ([]byte, string, error) {
	num := ""
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		num += n.String()
	}

	c := gg.NewContext(300, 120)
	if err := c.LoadFontFaceFromBytes(techno, 60); err != nil {
		panic(err)
	}
	c.MeasureString(num)
	c.SetHexColor("#FFFFFF")
	c.Clear()

	c.SetRGB(0, 0, 0)
	c.Fill()

	c.DrawStringWrapped(num, 20, 30, 0, 0, 300, 1.5, gg.AlignLeft)
	//c.SetRGB(0, 0, 0)
	//c.Fill()
	buf := new(bytes.Buffer)
	err := png.Encode(buf, c.Image())
	if err != nil {
		return nil, num, err
	}
	return buf.Bytes(), num, nil
}
