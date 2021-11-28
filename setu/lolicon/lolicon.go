package lolicon

import (
	"OPQBot-QQGroupManager/Core"
	"OPQBot-QQGroupManager/setu/setucore"
	"github.com/mcoo/requests"
	"github.com/sirupsen/logrus"
	"strings"
)

var log *logrus.Entry

type Provider struct {
}

func (p *Provider) InitProvider(l *logrus.Entry, bot *Core.Bot) {
	log = l
	log.Info("启动成功")
}
func (p *Provider) SearchPic(word string, r18 bool, num int) ([]setucore.Pic, error) {
	return nil, nil
}
func (p *Provider) SearchPicFromUser(word string, userId string, r18 bool, num int) ([]setucore.Pic, error) {
	return nil, nil
}
func (p *Provider) SearchPicToDB() (num int, e error) {
	res, err := requests.Get("https://api.lolicon.app/setu/v2?num=100&proxy=false&r18=2")
	if err != nil {
		return 0, err
	}
	var setus SetuRes
	if err = res.Json(&setus); err != nil {
		return 0, err
	}
	for _, v := range setus.Data {
		pic := setucore.Pic{
			Id:             v.Pid,
			Page:           v.P,
			Title:          v.Title,
			Author:         v.Author,
			AuthorID:       v.Uid,
			OriginalPicUrl: v.Urls.Original,
			Tag:            strings.Join(v.Tags, ","),
			R18:            v.R18,
			LastSendTime:   0,
		}
		err = setucore.AddPicToDB(pic)
		if err == nil {
			num += 1
		}
	}
	return
}

type SetuRes struct {
	Error string `json:"error"`
	Data  []struct {
		Pid        int      `json:"pid"`
		P          int      `json:"p"`
		Uid        int      `json:"uid"`
		Title      string   `json:"title"`
		Author     string   `json:"author"`
		R18        bool     `json:"r18"`
		Width      int      `json:"width"`
		Height     int      `json:"height"`
		Tags       []string `json:"tags"`
		Ext        string   `json:"ext"`
		UploadDate int64    `json:"uploadDate"`
		Urls       struct {
			Original string `json:"original"`
		} `json:"urls"`
	} `json:"data"`
}
