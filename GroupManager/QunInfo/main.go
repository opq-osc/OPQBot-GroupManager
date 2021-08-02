package QunInfo

import (
	"OPQBot-QQGroupManager/Core"
	"errors"
	"github.com/mcoo/OPQBot/qzone"
	"github.com/mcoo/requests"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type Qun struct {
	b          *Core.Bot
	RetryCount int
	QQ         string
	Gtk        string
	Gtk2       string
	PSkey      string
	Skey       string
	Uin        string
}

func NewQun(b *Core.Bot) Qun {
	q := Qun{b: b, RetryCount: 5}
	q.GetCookie()
	return q
}
func (q *Qun) GetBkn() string {
	return qzone.GenderGTK(q.Skey)
}
func (q *Qun) GetReq() *requests.Request {
	r := requests.Requests()
	c := &http.Cookie{
		Name:  "pt2gguin",
		Value: q.Uin,
	}
	r.SetCookie(c)
	c = &http.Cookie{
		Name:  "uin",
		Value: q.Uin,
	}
	r.SetCookie(c)
	c = &http.Cookie{
		Name:  "skey",
		Value: q.Skey,
	}
	r.SetCookie(c)
	c = &http.Cookie{
		Name:  "p_skey",
		Value: q.PSkey,
	}
	r.SetCookie(c)
	c = &http.Cookie{
		Name:  "p_uin",
		Value: q.Uin,
	}
	r.SetCookie(c)
	log.Println(r.Cookies)
	return r
}
func (q *Qun) GetCookie() {
	cookie, _ := q.b.GetUserCookie()
	q.Skey = cookie.Skey
	q.PSkey = cookie.PSkey.Qun
	q.Gtk = qzone.GenderGTK(cookie.Skey)
	q.Gtk2 = qzone.GenderGTK(cookie.PSkey.Qun)
	q.QQ = strconv.FormatInt(q.b.QQ, 10)
	q.Uin = "o" + q.QQ
}
func (q *Qun) GetInfo() {
	req := q.GetReq()
	res, _ := req.Get("https://qun.qq.com/cgi-bin/qunwelcome/myinfo?callback=?&bkn=" + q.GetBkn())
	log.Println(res.Text())
}

type GroupInfoResult struct {
	Retcode int    `json:"retcode"`
	Msg     string `json:"msg"`
	Data    struct {
		GroupInfo struct {
			GroupCode   string `json:"groupCode"`
			GroupName   string `json:"groupName"`
			GroupMember int    `json:"groupMember"`
			CreateDate  string `json:"createDate"`
		} `json:"groupInfo"`
		ActiveData struct {
			ActiveData  int `json:"activeData"`
			GroupMember int `json:"groupMember"`
			Ratio       int `json:"ratio"`
			DataList    []struct {
				Date   int `json:"date"`
				Number int `json:"number"`
			} `json:"dataList"`
		} `json:"activeData"`
		MsgInfo struct {
			Total    int `json:"total"`
			DataList []struct {
				Date   int `json:"date"`
				Number int `json:"number"`
			} `json:"dataList"`
		} `json:"msgInfo"`
		JoinData struct {
			Total    int `json:"total"`
			DataList []struct {
				Date   int `json:"date"`
				Number int `json:"number"`
			} `json:"dataList"`
		} `json:"joinData"`
		ExitData struct {
			Total    int `json:"total"`
			DataList []struct {
				Date   int `json:"date"`
				Number int `json:"number"`
			} `json:"dataList"`
		} `json:"exitData"`
		ApplyData struct {
			Total    int `json:"total"`
			DataList []struct {
				Date   int `json:"date"`
				Number int `json:"number"`
			} `json:"dataList"`
		} `json:"applyData"`
		MemberData struct {
			Total    int `json:"total"`
			DataList []struct {
				Date   int `json:"date"`
				Number int `json:"number"`
			} `json:"dataList"`
		} `json:"memberData"`
		LastDataTime int `json:"lastDataTime"`
	} `json:"data"`
}
type GroupMembersResult struct {
	Retcode int    `json:"retcode"`
	Msg     string `json:"msg"`
	Data    struct {
		ListNext  int `json:"listNext"`
		SpeakRank []struct {
			Uin      string `json:"uin"`
			Avatar   string `json:"avatar"`
			Nickname string `json:"nickname"`
			Active   int    `json:"active"`
			MsgCount int    `json:"msgCount"`
		} `json:"speakRank"`
	} `json:"data"`
}

func (q *Qun) GetGroupInfo(groupId int64, time int) (result GroupInfoResult, e error) {
	if q.RetryCount <= 0 {
		return result, errors.New("超过重试次数")
	}
	req := q.GetReq()
	strGroupId := strconv.FormatInt(groupId, 10)
	req.Header.Set("qname-service", "976321:131072")
	req.Header.Set("qname-space", "Production")
	req.Header.Set("referer", "https://qun.qq.com/m/qun/activedata/active.html?_wv=3&_wwv=128&gc="+strGroupId+"&src=2")

	res, e := req.Get("https://qun.qq.com/m/qun/activedata/proxy/domain/qun.qq.com/cgi-bin/manager/report/index?gc=" + strGroupId + "&time=" + strconv.Itoa(time) + "&bkn=" + q.GetBkn())
	if e != nil {
		return result, e
	}
	e = res.Json(&result)
	if e != nil {

		log.Println(res.R.Request.RequestURI, res.Text())
		return result, e
	}
	if result.Retcode == 100000 {
		q.RetryCount -= 1
		return q.GetGroupInfo(groupId, time)
	}
	q.RetryCount = 5
	return result, nil
}
func (q *Qun) GetGroupMembersInfo(groupId int64, time int) (result GroupMembersResult, e error) {
	if q.RetryCount <= 0 {
		return result, errors.New("超过重试次数")
	}
	req := q.GetReq()
	strGroupId := strconv.FormatInt(groupId, 10)
	req.Header.Set("qname-service", "976321:131072")
	req.Header.Set("qname-space", "Production")
	req.Header.Set("referer", "https://qun.qq.com/m/qun/activedata/active.html?_wv=3&_wwv=128&gc="+strGroupId+"&src=2")
	res, e := req.Get("https://qun.qq.com/m/qun/activedata/proxy/domain/qun.qq.com/cgi-bin/manager/report/list?gc=" + strGroupId + "&time=" + strconv.Itoa(time) + "&bkn=" + q.GetBkn() + "&type=0&start=0")
	if e != nil {
		return result, e
	}
	e = res.Json(&result)
	if e != nil {

		log.Println(res.R.Request.RequestURI, res.Text())
		return result, e
	}
	if result.Retcode == 100000 {
		q.RetryCount -= 1
		return q.GetGroupMembersInfo(groupId, time)
	}
	sort.Slice(result.Data.SpeakRank, func(i, j int) bool {
		return result.Data.SpeakRank[i].Active > result.Data.SpeakRank[j].Active
	})
	q.RetryCount = 5
	return result, nil
}
