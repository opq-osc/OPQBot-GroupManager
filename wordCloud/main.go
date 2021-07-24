package wordCloud

import (
	"OPQBot-QQGroupManager/Core"
	"bytes"
	"github.com/huichen/sego"
	"github.com/mcoo/OPQBot"
	"github.com/psykhi/wordclouds"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"image/color"
	"image/png"
	"strings"
	"time"
)

var DefaultColors = []color.RGBA{
	{0x1b, 0x1b, 0x1b, 0xff},
	{0x48, 0x48, 0x4B, 0xff},
	{0x59, 0x3a, 0xee, 0xff},
	{0x65, 0xCD, 0xFA, 0xff},
	{0x70, 0xD6, 0xBF, 0xff},
}

type Module struct {
	db         *gorm.DB
	MsgChannel chan OPQBot.GroupMsgPack
}

var (
	log *logrus.Entry
)

func (m *Module) ModuleInfo() Core.ModuleInfo {
	return Core.ModuleInfo{
		Name:        "词云生成",
		Author:      "enjoy",
		Description: "给群生成聊天词云",
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
	var segmented sego.Segmenter
	segmented.LoadDictionary("./dictionary.txt")
	for {
		msg := <-m.MsgChannel
		split := strings.Split(sego.SegmentsToString(segmented.Segment([]byte(msg.Content)), false), " ")
		for _, v := range split {
			if s := strings.Split(v, "/"); len(s) == 2 && (len(s[0]) > 1) {
				err := m.AddHotWord(s[0], msg.FromGroupID)
				if err != nil {
					log.Error(err)
				}
			}
		}

	}
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
	err = b.AddEvent(OPQBot.EventNameOnGroupMessage, func(qq int64, packet *OPQBot.GroupMsgPack) {
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
				sendMsg := "今日本群词云\n"
				hotMap := map[string]int{}
				for i := 0; i < len(hotWords); i++ {
					hotMap[hotWords[i].Word] = hotWords[i].Count
				}
				log.Info(hotMap)
				colors := make([]color.Color, 0)
				for _, c := range DefaultColors {
					colors = append(colors, c)
				}

				img := wordclouds.NewWordcloud(hotMap, wordclouds.FontMaxSize(128), wordclouds.FontMinSize(60), wordclouds.FontFile("./font.ttf"),
					wordclouds.Height(1024),
					wordclouds.Width(2048), wordclouds.Colors(colors)).Draw()

				buf := new(bytes.Buffer)
				err = png.Encode(buf, img)
				if err != nil {
					log.Error(err)
					return
				}
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
	err := m.db.Debug().Where(" (? - hot_time) <= 86400 AND group_id = ? AND word = ?", t.Unix(), groupId, word).Find(&hotWord).Error
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
	err := m.db.Where("(? - hot_time) <= 86400  AND group_id = ?", t, groupId).Find(&hotWord).Error
	if err != nil {
		return nil, err
	}

	return hotWord, nil
}
func (m *Module) GetWeeklyWord(groupId int64) ([]HotWord, error) {
	var hotWord []HotWord
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Add(24*time.Hour).Format("2006-01-02")+" 00:00:00", time.Local)

	err := m.db.Where("(? - hot_time) <= 604800 AND group_id = ?", t, groupId).Find(&hotWord).Error
	if err != nil {
		return nil, err
	}

	return hotWord, nil
}
func init() {
	Core.RegisterModule(&Module{})
}
