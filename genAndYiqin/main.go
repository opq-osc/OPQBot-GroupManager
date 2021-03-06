package genAndYiqin

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/Core"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/mcoo/OPQBot"

	"io/ioutil"

	"net/http"

	"strings"
)

type GenChaxunRes struct {
	Data []DataRes `json:"data"`
}

type DataRes struct {
	Definitions []Items    `json:"definitions"`
	Tags        []TageItem `json:"tags"`
}

type Items struct {
	Content   string      `json:"content"`
	Plaintext string      `json:"plaintext"`
	Images    []ImageItem `json:"images"`
}

type ImageItem struct {
	Full FullItem `json:"full"`
}

type FullItem struct {
	Path string `json:"path"`
}
type TageItem struct {
	Name string `json:"name"`
}
type Module struct {
}
type YiqingRes struct {
	Title      string    `json:"title"`
	Time       string    `json:"time"`
	IncrTime   string    `json:"incrTime"`
	logcation  Logcation `json:"logcation"`
	Colums     []Colums  `json:"colums"`
	MainReport struct {
		Id       int    `json:"id"`
		Area     string `json:"area"`
		Report   string `json:"report"`
		Dateline string `json:"dateline"`
		Date     int64  `json:"date"`
	} `json:"mainReport"`
	ContryData struct {
		SureCnt         string `json:"sure_cnt"`
		SureNewCnt      string `json:"sure_new_cnt"`
		RestSureCnt     string `json:"rest_sure_cnt"`
		RestSureCntIncr string `json:"rest_sure_cnt_incr"`
		InputCnt        string `json:"input_cnt"`
		HiddenCnt       string `json:"hidden_cnt"`
		HiddenCntIncr   string `json:"hidden_cnt_incr"`
		CureCnt         string `json:"cure_cnt"`
		YstCureCnt      string `json:"yst_cure_cnt"`
		YstDieCnt       string `json:"yst_die_cnt"`
		YstLikeCnt      string `json:"yst_like_cnt"`
		YstSureCnt      string `json:"yst_sure_cnt"`
		YstSureHid      string `json:"yst_sure_hid"`
	}
}

type Colums struct {
	Title string `json:"title"`
	List  []List `json:"list"`
}

type List struct {
	Current int64  `json:"current"`
	Incr    string `json:"incr"`
}

type Logcation struct {
	Province string `json:"province"`
	City     string `json:"city"`
}

var log *logrus.Entry

func (m *Module) ModuleInfo() Core.ModuleInfo {
	return Core.ModuleInfo{
		Name:        "????????????????????????",
		Author:      "bypanghu",
		Description: "",
		Version:     0,
	}
}
func (m *Module) ModuleInit(b *Core.Bot, l *logrus.Entry) error {
	log = l
	b.BotCronManager.AddJob(-1, "Yiqing", "0 8,18 * * *", func() {
		client := &http.Client{}
		baseUrl := "https://m.sm.cn/api/rest?method=Huoshenshan.local"
		req, err := http.NewRequest("GET", baseUrl, nil)
		req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
		req.Header.Add("referer", "https://broccoli.uc.cn/")
		if err != nil {
			panic(err)
		}
		response, _ := client.Do(req)
		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		var res YiqingRes
		json.Unmarshal(s, &res)
		ups := fmt.Sprintf("????????????")
		ups += fmt.Sprintf("%s-%s\n??????????????????%s\n", res.Title, res.Time, res.MainReport.Report)
		ups += fmt.Sprintf("[??????190][??????190][??????190]????????????[??????190][??????190][??????190]\n")
		ups += fmt.Sprintf("[??????145]??????????????????%s???????????????%s???\n", res.ContryData.SureCnt, res.ContryData.YstCureCnt)
		ups += fmt.Sprintf("[??????145]??????????????????%s???????????????%s???\n", res.ContryData.RestSureCnt, res.ContryData.RestSureCntIncr)
		ups += fmt.Sprintf("[??????145]??????????????????%s???\n", res.ContryData.InputCnt)
		ups += fmt.Sprintf("[??????145]??????????????????%s???????????????%s???\n", res.ContryData.CureCnt, res.ContryData.YstCureCnt)
		ups += fmt.Sprintf("[??????66][??????66][??????66]????????????????????????????????????")
		b.SendGroupTextMsg(-1, fmt.Sprintf(ups))
		fmt.Println(ups)
	})

	_, err := b.AddEvent(OPQBot.EventNameOnGroupMessage, func(botQQ int64, packet *OPQBot.GroupMsgPack) {
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
		if len(cm) == 2 && cm[0] == "?????????" {
			b.SendGroupTextMsg(packet.FromGroupID, fmt.Sprintf("???????????????%s", cm[1]))
			client := &http.Client{}
			baseUrl := "https://api.jikipedia.com/go/search_entities"
			postData := make(map[string]interface{})
			postData["phrase"] = cm[1]
			postData["page"] = 1
			bytesData, err := json.Marshal(postData)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			reader := bytes.NewReader(bytesData)
			req, err := http.NewRequest("POST", baseUrl, reader)
			req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
			req.Header.Add("referer", "https://broccoli.uc.cn/")
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			if err != nil {
				panic(err)
			}
			response, _ := client.Do(req)
			defer response.Body.Close()
			s, err := ioutil.ReadAll(response.Body)
			var res GenChaxunRes
			json.Unmarshal(s, &res)
			var content string
			for i, a := range res.Data {
				if i == 1 {
					for j, b := range a.Definitions {
						if j == 0 {
							content = b.Plaintext
						}
					}
				}
			}
			if content == "" {
				b.SendGroupTextMsg(packet.FromGroupID, "?????????????????????")
			} else {
				b.SendGroupTextMsg(packet.FromGroupID, fmt.Sprintf("%s", content))
			}
			return
		}
		if packet.Content == "????????????" {
			b.SendGroupTextMsg(packet.FromGroupID, "??????????????????")
			client := &http.Client{}
			baseUrl := "https://m.sm.cn/api/rest?method=Huoshenshan.local"
			req, err := http.NewRequest("GET", baseUrl, nil)
			req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
			req.Header.Add("referer", "https://broccoli.uc.cn/")
			if err != nil {
				panic(err)
			}
			response, _ := client.Do(req)
			defer response.Body.Close()
			s, err := ioutil.ReadAll(response.Body)
			var res YiqingRes
			json.Unmarshal(s, &res)
			ups := fmt.Sprintf("%s-%s\n??????????????????%s\n", res.Title, res.Time, res.MainReport.Report)
			ups += fmt.Sprintf("[??????190][??????190][??????190]????????????[??????190][??????190][??????190]\n")
			ups += fmt.Sprintf("[??????145]??????????????????%s???????????????%s???\n", res.ContryData.SureCnt, res.ContryData.YstCureCnt)
			ups += fmt.Sprintf("[??????145]??????????????????%s???????????????%s???\n", res.ContryData.RestSureCnt, res.ContryData.RestSureCntIncr)
			ups += fmt.Sprintf("[??????145]??????????????????%s???\n", res.ContryData.InputCnt)
			ups += fmt.Sprintf("[??????145]??????????????????%s???????????????%s???\n", res.ContryData.CureCnt, res.ContryData.YstCureCnt)
			ups += fmt.Sprintf("[??????66][??????66][??????66]????????????????????????????????????")
			b.SendGroupTextMsg(packet.FromGroupID, fmt.Sprintf(ups))
			log.Println(ups)
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
