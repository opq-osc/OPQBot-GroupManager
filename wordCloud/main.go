package wordCloud

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/Core"
	"bytes"
	"errors"
	"fmt"
	"github.com/go-ego/gse"
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/requests"
	"github.com/mcoo/wordclouds"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"image/color"
	"image/png"
	"time"
)

var DefaultColors = []color.RGBA{
	{0x1b, 0x1b, 0x1b, 0xff},
	{0x48, 0x48, 0x4B, 0xff},
	{0x59, 0x3a, 0xee, 0xff},
	{0x65, 0xCD, 0xFA, 0xff},
	{0x70, 0xD6, 0xBF, 0xff},
	{153, 50, 204, 255},
	{100, 149, 237, 255},
	{0, 255, 127, 255},
	{255, 0, 0, 255},
}

type Module struct {
	db         *gorm.DB
	MsgChannel chan OPQBot.GroupMsgPack
	ImgServer  string
}

var (
	log *logrus.Entry
)

func (m *Module) ModuleInfo() Core.ModuleInfo {
	return Core.ModuleInfo{
		Name:        "词云生成",
		Author:      "enjoy",
		Description: "给群生成聊天词云 还可以查询奥运信息呢！",
		Version:     0,
	}
}

type HotWord struct {
	gorm.Model
	GroupId int64  `gorm:"index"`
	Word    string `gorm:"index"`
	Count   int    `gorm:"not null"`
	HotTime int64  `gorm:"index;not null"`
}

func (m *Module) DoHotWord() {
	var segmenter gse.Segmenter
	err := segmenter.LoadDict("./dictionary.txt")
	if err != nil {
		log.Error(err)
	}
	//var segmented sego.Segmenter
	//segmented.LoadDictionary("./dictionary.txt")
	for {
		msg := <-m.MsgChannel
		tmp := segmenter.Segment([]byte(msg.Content))
		for _, v := range gse.ToSlice(tmp, false) {
			if len([]rune(v)) > 1 {
				err := m.AddHotWord(v, msg.FromGroupID)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Data string `json:"data"`
	} `json:"result"`
}
type ReqStruct struct {
	To      string `json:"to"`
	Fetcher struct {
		Name   string `json:"name"`
		Params struct {
			Data string `json:"data"`
		} `json:"params"`
	} `json:"fetcher"`
	Converter struct {
		Extend struct {
			JavascriptDelay string `json:"javascript-delay"`
		} `json:"extend"`
	} `json:"converter"`
	Template string `json:"template"`
}

func (m *Module) GetUrlPic(url string, width, height, javascriptDelay int) (string, error) {
	r, err := requests.PostJson(m.ImgServer, fmt.Sprintf(`{
  "to": "image",
  "converter": {
    "uri": "%s",
    "width": %d,
    "height": %d,
	"extend": {
          "javascript-delay": "%d"
	}
  },
  "template": ""
}`, url, width, height, javascriptDelay))
	if err != nil {
		return "", err
	}
	var result Result
	err = r.Json(&result)
	if err != nil {
		return "", err
	}
	if result.Code != 0 {
		log.Error()
		return "", errors.New(fmt.Sprintf("[%d] %s", result.Code, result.Message))
	}
	return result.Result.Data, nil
}
func (m *Module) ModuleInit(b *Core.Bot, l *logrus.Entry) error {
	log = l
	m.db = b.DB
	m.MsgChannel = make(chan OPQBot.GroupMsgPack, 30)
	go m.DoHotWord()
	err := b.DB.AutoMigrate(&HotWord{})
	if err != nil {
		return err
	}
	Config.Lock.RLock()
	m.ImgServer = Config.CoreConfig.HtmlToImgUrl
	Config.Lock.RUnlock()

	_, err = b.AddEvent(OPQBot.EventNameOnGroupMessage, func(qq int64, packet *OPQBot.GroupMsgPack) {
		if packet.FromUserID != b.QQ {
			if packet.MsgType == "TextMsg" {
				m.MsgChannel <- *packet
			}
			if packet.Content == "今日词云" {
				hotWords, err := m.GetTodayWord(packet.FromGroupID)
				if err != nil {
					log.Error(err)
					return
				}
				sendMsg := "今日本群词云"
				hotMap := map[string]int{}
				for i := 0; i < len(hotWords); i++ {
					if len([]rune(hotWords[i].Word)) <= 1 {
						continue
					}
					hotMap[hotWords[i].Word] = hotWords[i].Count
				}
				log.Info(hotMap)
				colors := make([]color.Color, 0)
				for _, c := range DefaultColors {
					colors = append(colors, c)
				}
				b.PrintMemStats()
				img := wordclouds.NewWordcloud(hotMap, wordclouds.FontMaxSize(200), wordclouds.FontMinSize(40), wordclouds.FontFile("./font.ttf"),
					wordclouds.Height(1324),
					wordclouds.Width(1324), wordclouds.Colors(colors)).Draw()
				b.PrintMemStats()
				buf := new(bytes.Buffer)
				err = png.Encode(buf, img)
				if err != nil {
					log.Error(err)
					return
				}
				b.PrintMemStats()
				b.SendGroupPicMsg(packet.FromGroupID, sendMsg, buf.Bytes())
			}
			if packet.Content == "奥运赛程" {
				pic, err := m.GetUrlPic("https://tiyu.baidu.com/tokyoly/delegation/8567/tab/%E8%B5%9B%E7%A8%8B/type/all", 360, 640, 0)
				if err != nil {
					log.Error(err)
					return
				}
				b.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content: OPQBot.SendTypePicMsgByBase64Content{
						Content: "",
						Base64:  pic,
						Flash:   false,
					},
					CallbackFunc: nil,
				})
				return
			}
			if packet.Content == "中国奥运" {
				pic, err := m.GetUrlPic("https://tiyu.baidu.com/tokyoly/delegation/8567/tab/%E5%A5%96%E7%89%8C%E6%98%8E%E7%BB%86", 360, 640, 0)
				if err != nil {
					log.Error(err)
					return
				}
				b.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content: OPQBot.SendTypePicMsgByBase64Content{
						Content: "",
						Base64:  pic,
						Flash:   false,
					},
					CallbackFunc: nil,
				})
				return
			}
			if packet.Content == "奥运" {
				pic, err := m.GetUrlPic("https://tiyu.baidu.com/tokyoly/home/tab/%E5%A5%96%E7%89%8C%E6%A6%9C", 360, 640, 0)
				if err != nil {
					log.Error(err)
					return
				}
				b.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content: OPQBot.SendTypePicMsgByBase64Content{
						Content: "",
						Base64:  pic,
						Flash:   false,
					},
					CallbackFunc: nil,
				})
				return
			}
			if packet.Content == "本周词云" {
				hotWords, err := m.GetWeeklyWord(packet.FromGroupID)
				if err != nil {
					log.Error(err)
					return
				}
				sendMsg := "本周词云"
				hotMap := map[string]int{}
				for i := 0; i < len(hotWords); i++ {
					if len([]rune(hotWords[i].Word)) <= 1 {
						continue
					}
					hotMap[hotWords[i].Word] = hotWords[i].Count
				}
				log.Info(hotMap)
				colors := make([]color.Color, 0)
				for _, c := range DefaultColors {
					colors = append(colors, c)
				}

				b.PrintMemStats()
				img := wordclouds.NewWordcloud(hotMap, wordclouds.FontMaxSize(200), wordclouds.FontMinSize(40), wordclouds.FontFile("./font.ttf"),
					wordclouds.Height(1324),
					wordclouds.Width(1324), wordclouds.Colors(colors)).Draw()
				b.PrintMemStats()
				buf := new(bytes.Buffer)
				err = png.Encode(buf, img)
				if err != nil {
					log.Error(err)
					return
				}
				b.PrintMemStats()
				b.SendGroupPicMsg(packet.FromGroupID, sendMsg, buf.Bytes())
			}
		}
	})
	if err != nil {
		return err
	}
	return nil
}
func (m *Module) AddHotWord(word string, groupId int64) error {
	var hotWord []HotWord

	t, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Add(24*time.Hour).Format("2006-01-02")+" 00:00:00", time.Local)
	err := m.db.Where(" (? - hot_time) <= 86400 AND group_id = ? AND word = ?", t.Unix(), groupId, word).Find(&hotWord).Error
	if err != nil {
		return err
	}

	if len(hotWord) > 0 {
		err := m.db.Model(&hotWord[0]).Update("count", hotWord[0].Count+1).Error
		if err != nil {
			return err
		}
	} else {
		tmp := HotWord{
			GroupId: groupId,
			Word:    word,
			Count:   1,
			HotTime: time.Now().Unix(),
		}
		err := m.db.Create(&tmp).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func (m *Module) GetTodayWord(groupId int64) ([]HotWord, error) {
	var hotWord []HotWord
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Add(24*time.Hour).Format("2006-01-02")+" 00:00:00", time.Local)
	err := m.db.Where("(? - hot_time) <= 86400  AND group_id = ?", t, groupId).Limit(100).Find(&hotWord).Error
	if err != nil {
		return nil, err
	}

	return hotWord, nil
}
func (m *Module) GetWeeklyWord(groupId int64) ([]HotWord, error) {
	var hotWord []HotWord
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Add(24*time.Hour).Format("2006-01-02")+" 00:00:00", time.Local)

	err := m.db.Where("(? - hot_time) <= 604800 AND group_id = ?", t, groupId).Limit(100).Find(&hotWord).Error
	if err != nil {
		return nil, err
	}

	return hotWord, nil
}
func init() {
	Core.RegisterModule(&Module{})
}
