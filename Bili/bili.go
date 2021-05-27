package bili

import (
	"OPQBot-QQGroupManager/Config"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/mcoo/requests"
)

type Manager struct {
	ups        map[int64]Up
	upsMapLock *sync.RWMutex
	r          *requests.Request
}
type SearchResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Seid           string `json:"seid"`
		Page           int    `json:"page"`
		Pagesize       int    `json:"pagesize"`
		Numresults     int    `json:"numResults"`
		Numpages       int    `json:"numPages"`
		SuggestKeyword string `json:"suggest_keyword"`
		RqtType        string `json:"rqt_type"`
		CostTime       struct {
			ParamsCheck         string `json:"params_check"`
			GetUpuserLiveStatus string `json:"get upuser live status"`
			IllegalHandler      string `json:"illegal_handler"`
			AsResponseFormat    string `json:"as_response_format"`
			AsRequest           string `json:"as_request"`
			SaveCache           string `json:"save_cache"`
			DeserializeResponse string `json:"deserialize_response"`
			AsRequestFormat     string `json:"as_request_format"`
			Total               string `json:"total"`
			MainHandler         string `json:"main_handler"`
		} `json:"cost_time"`
		ExpList struct {
			Num5502 bool `json:"5502"`
			Num6600 bool `json:"6600"`
		} `json:"exp_list"`
		EggHit int `json:"egg_hit"`
		Result []struct {
			Type       string `json:"type"`
			Mid        int64  `json:"mid"`
			Uname      string `json:"uname"`
			Usign      string `json:"usign"`
			Fans       int    `json:"fans"`
			Videos     int    `json:"videos"`
			Upic       string `json:"upic"`
			VerifyInfo string `json:"verify_info"`
			Level      int    `json:"level"`
			Gender     int    `json:"gender"`
			IsUpuser   int    `json:"is_upuser"`
			IsLive     int    `json:"is_live"`
			RoomID     int    `json:"room_id"`
			Res        []struct {
				Aid          int    `json:"aid"`
				Bvid         string `json:"bvid"`
				Title        string `json:"title"`
				Pubdate      int    `json:"pubdate"`
				Arcurl       string `json:"arcurl"`
				Pic          string `json:"pic"`
				Play         string `json:"play"`
				Dm           int    `json:"dm"`
				Coin         int    `json:"coin"`
				Fav          int    `json:"fav"`
				Desc         string `json:"desc"`
				Duration     string `json:"duration"`
				IsPay        int    `json:"is_pay"`
				IsUnionVideo int    `json:"is_union_video"`
			} `json:"res"`
			OfficialVerify struct {
				Type int    `json:"type"`
				Desc string `json:"desc"`
			} `json:"official_verify"`
			HitColumns []interface{} `json:"hit_columns"`
		} `json:"result"`
		ShowColumn int `json:"show_column"`
	} `json:"data"`
}
type Up struct {
	Name    string
	Created int64
	Groups  []int64
}

type GroupNode struct {
	GroupId int64
	Next    *GroupNode
}

type BiliResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		List struct {
			Vlist []Video `json:"vlist"`
		} `json:"list"`
		Page struct {
			Pn    int `json:"pn"`
			Ps    int `json:"ps"`
			Count int `json:"count"`
		} `json:"page"`
	} `json:"data"`
}

type Video struct {
	Comment        int    `json:"comment"`
	Typeid         int    `json:"typeid"`
	Play           int    `json:"play"`
	Pic            string `json:"pic"`
	Subtitle       string `json:"subtitle"`
	Description    string `json:"description"`
	Copyright      string `json:"copyright"`
	Title          string `json:"title"`
	Review         int    `json:"review"`
	Author         string `json:"author"`
	Mid            int64  `json:"mid"`
	Created        int64  `json:"created"`
	Length         string `json:"length"`
	VideoReview    int    `json:"video_review"`
	Aid            int64  `json:"aid"`
	Bvid           string `json:"bvid"`
	HideClick      bool   `json:"hide_click"`
	IsPay          int    `json:"is_pay"`
	IsUnionVideo   int    `json:"is_union_video"`
	IsSteinsGate   int    `json:"is_steins_gate"`
	IsLivePlayback int    `json:"is_live_playback"`
}

type UpInfoResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    UpInfo `json:"data"`
}

type UpInfo struct {
	Card struct {
		Mid         string        `json:"mid"`
		Name        string        `json:"name"`
		Approve     bool          `json:"approve"`
		Sex         string        `json:"sex"`
		Rank        string        `json:"rank"`
		Face        string        `json:"face"`
		Displayrank string        `json:"DisplayRank"`
		Regtime     int           `json:"regtime"`
		Spacesta    int           `json:"spacesta"`
		Birthday    string        `json:"birthday"`
		Place       string        `json:"place"`
		Description string        `json:"description"`
		Article     int           `json:"article"`
		Attentions  []interface{} `json:"attentions"`
		Fans        int           `json:"fans"`
		Friend      int           `json:"friend"`
		Attention   int           `json:"attention"`
		Sign        string        `json:"sign"`
		LevelInfo   struct {
			CurrentLevel int `json:"current_level"`
			CurrentMin   int `json:"current_min"`
			CurrentExp   int `json:"current_exp"`
			NextExp      int `json:"next_exp"`
		} `json:"level_info"`
		Pendant struct {
			Pid               int    `json:"pid"`
			Name              string `json:"name"`
			Image             string `json:"image"`
			Expire            int    `json:"expire"`
			ImageEnhance      string `json:"image_enhance"`
			ImageEnhanceFrame string `json:"image_enhance_frame"`
		} `json:"pendant"`
		Nameplate struct {
			Nid        int    `json:"nid"`
			Name       string `json:"name"`
			Image      string `json:"image"`
			ImageSmall string `json:"image_small"`
			Level      string `json:"level"`
			Condition  string `json:"condition"`
		} `json:"nameplate"`
		Official struct {
			Role  int    `json:"role"`
			Title string `json:"title"`
			Desc  string `json:"desc"`
			Type  int    `json:"type"`
		} `json:"Official"`
		OfficialVerify struct {
			Type int    `json:"type"`
			Desc string `json:"desc"`
		} `json:"official_verify"`
		Vip struct {
			Type       int   `json:"type"`
			Status     int   `json:"status"`
			DueDate    int64 `json:"due_date"`
			VipPayType int   `json:"vip_pay_type"`
			ThemeType  int   `json:"theme_type"`
			Label      struct {
				Path        string `json:"path"`
				Text        string `json:"text"`
				LabelTheme  string `json:"label_theme"`
				TextColor   string `json:"text_color"`
				BgStyle     int    `json:"bg_style"`
				BgColor     string `json:"bg_color"`
				BorderColor string `json:"border_color"`
			} `json:"label"`
			AvatarSubscript    int    `json:"avatar_subscript"`
			NicknameColor      string `json:"nickname_color"`
			Role               int    `json:"role"`
			AvatarSubscriptURL string `json:"avatar_subscript_url"`
			Viptype            int    `json:"vipType"`
			Vipstatus          int    `json:"vipStatus"`
		} `json:"vip"`
	} `json:"card"`
	Following    bool `json:"following"`
	ArchiveCount int  `json:"archive_count"`
	ArticleCount int  `json:"article_count"`
	Follower     int  `json:"follower"`
}

func NewManager() (m Manager) {
	Config.Lock.RLock()
	defer Config.Lock.RUnlock()
	m.ups = make(map[int64]Up)
	for groupId, v := range Config.CoreConfig.GroupConfig {
		for mid, v1 := range v.BiliUps {
			if up, ok := m.ups[mid]; ok {
				up.Groups = append(up.Groups, groupId)
			} else {
				m.ups[mid] = Up{
					Name:    v1.Name,
					Created: v1.Created,
					Groups:  []int64{groupId},
				}
			}
		}
	}

	m.upsMapLock = &sync.RWMutex{}
	m.r = requests.Requests()
	m.r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36 Edg/90.0.818.66")
	return
}

func (m *Manager) GetAllSubscribeUp() map[int64]Up {
	m.upsMapLock.RLock()
	defer m.upsMapLock.RUnlock()
	return m.ups
}

func (m *Manager) ScanUpdate() (result []Video) {
	m.upsMapLock.Lock()
	defer m.upsMapLock.Unlock()
	Config.Lock.Lock()
	defer Config.Lock.Unlock()
	for mid, v := range m.ups {
		res, err := m.r.Get(fmt.Sprintf("https://api.bilibili.com/x/space/arc/search?mid=%d&ps=1", mid))
		if err != nil {
			log.Println(err)
			break
		}
		var biliResult BiliResult
		err = res.Json(&biliResult)
		if err != nil {
			log.Println(err)
			break
		}
		if biliResult.Code == 0 {
			if len(biliResult.Data.List.Vlist) != 0 && biliResult.Data.List.Vlist[0].Created > v.Created {
				// 检测到更新
				if v.Created != 0 { // 跳过订阅后首次扫描
					result = append(result, biliResult.Data.List.Vlist[0])
				}
				v.Created = biliResult.Data.List.Vlist[0].Created
				m.ups[mid] = v

				for _, v1 := range v.Groups {
					if v2, ok := Config.CoreConfig.GroupConfig[v1]; ok {
						if v3, ok := v2.BiliUps[mid]; ok {
							v3.Created = v.Created
							v2.BiliUps[mid] = v3
						}
						Config.CoreConfig.GroupConfig[v1] = v2
					}
				}
				Config.Save()

			}
		}
	}
	return
}

func (m *Manager) GetUpInfoByMid(mid int64) (u UpInfoResult, e error) {
	res, e := m.r.Get(fmt.Sprintf("https://api.bilibili.com/x/web-interface/card?mid=%d", mid))
	if e != nil {
		return
	}
	e = res.Json(&u)
	if e != nil {
		return
	}
	return
}
func (m *Manager) UnSubscribeUp(groupId int64, mid int64) (e error) {
	m.upsMapLock.Lock()
	defer m.upsMapLock.Unlock()
	Config.Lock.Lock()
	defer Config.Lock.Unlock()
	if up, ok := m.ups[mid]; ok {
		for i, v := range up.Groups {
			if v == groupId {
				if len(up.Groups) == 1 {
					delete(m.ups, mid)
				} else {
					up.Groups = append(up.Groups[:i], up.Groups[i+1:]...)
					m.ups[mid] = up
				}

				if v2, ok := Config.CoreConfig.GroupConfig[groupId]; ok {
					delete(v2.BiliUps, mid)
					Config.CoreConfig.GroupConfig[groupId] = v2
				}
				Config.Save()

				return nil
			}
		}

		return nil
	} else {
		return nil
	}

}
func (m *Manager) SubscribeUpByKeyword(groupId int64, keyword string) (u UpInfoResult, e error) {
	if groupId == 0 {
		e = errors.New("默认群禁止订阅!")
		return
	}
	mid, e := m.SearchUp(keyword)
	if e != nil {
		return
	}
	m.upsMapLock.Lock()
	defer m.upsMapLock.Unlock()
	Config.Lock.Lock()
	defer Config.Lock.Unlock()

	if up, ok := m.ups[mid]; ok {
		in := false
		for _, v := range up.Groups {
			if v == groupId {
				in = true
				break
			}
		}
		if in {
			e = errors.New("该群已经订阅了该UP")
			return
		}
	}

	if u, e = m.GetUpInfoByMid(mid); e != nil {
		return
	} else {
		if u.Code != 0 {
			if u.Code == -404 {
				e = errors.New("找不到该UP")
				return
			}
			e = errors.New("Code Err")
			return
		}

		if up, c := m.ups[mid]; c {
			up.Groups = append(up.Groups, groupId)
			m.ups[mid] = up
			if v, ok := Config.CoreConfig.GroupConfig[groupId]; ok {
				if v.BiliUps == nil {
					v.BiliUps = map[int64]Config.Up{}
				}
				v.BiliUps[mid] = Config.Up{Name: u.Data.Card.Name, Created: up.Created}
				Config.CoreConfig.GroupConfig[groupId] = v
			} else {
				if v.BiliUps == nil {
					v.BiliUps = map[int64]Config.Up{}
				}
				v = Config.CoreConfig.DefaultGroupConfig
				v.BiliUps[mid] = Config.Up{Name: u.Data.Card.Name, Created: up.Created}
				Config.CoreConfig.GroupConfig[groupId] = v
			}
			Config.Save()
		} else {
			up = Up{Name: u.Data.Card.Name, Groups: []int64{groupId}, Created: 0}
			m.ups[mid] = up
			if v, ok := Config.CoreConfig.GroupConfig[groupId]; ok {
				if v.BiliUps == nil {
					v.BiliUps = map[int64]Config.Up{}
				}
				v.BiliUps[mid] = Config.Up{Name: u.Data.Card.Name, Created: 0}
				Config.CoreConfig.GroupConfig[groupId] = v
			} else {
				if v.BiliUps == nil {
					v.BiliUps = map[int64]Config.Up{}
				}
				v = Config.CoreConfig.DefaultGroupConfig
				v.BiliUps[mid] = Config.Up{Name: u.Data.Card.Name, Created: 0}
				Config.CoreConfig.GroupConfig[groupId] = v
			}
			Config.Save()
		}

	}
	return
}
func (m *Manager) SubscribeUpByMid(groupId int64, mid int64) (u UpInfoResult, e error) {
	if groupId == 0 {
		e = errors.New("默认群禁止订阅!")
		return
	}
	m.upsMapLock.Lock()
	defer m.upsMapLock.Unlock()
	Config.Lock.Lock()
	defer Config.Lock.Unlock()
	if up, ok := m.ups[mid]; ok {
		in := false
		for _, v := range up.Groups {
			if v == groupId {
				in = true
				break
			}
		}
		if in {
			e = errors.New("该群已经订阅了该UP")
			return
		}
	}

	if u, e = m.GetUpInfoByMid(mid); e != nil {
		return
	} else {
		if u.Code != 0 {
			if u.Code == -404 {
				e = errors.New("找不到该UP")
				return
			}
			e = errors.New("Code Err")
			return
		}

		if up, c := m.ups[mid]; c {
			up.Groups = append(up.Groups, groupId)
			m.ups[mid] = up
			if v, ok := Config.CoreConfig.GroupConfig[groupId]; ok {
				if v.BiliUps == nil {
					v.BiliUps = map[int64]Config.Up{}
				}
				v.BiliUps[mid] = Config.Up{Name: u.Data.Card.Name, Created: up.Created}
				Config.CoreConfig.GroupConfig[groupId] = v
			} else {
				if v.BiliUps == nil {
					v.BiliUps = map[int64]Config.Up{}
				}
				v = Config.CoreConfig.DefaultGroupConfig
				v.BiliUps[mid] = Config.Up{Name: u.Data.Card.Name, Created: up.Created}
				Config.CoreConfig.GroupConfig[groupId] = v
			}
			Config.Save()
		} else {
			up = Up{Name: u.Data.Card.Name, Groups: []int64{groupId}, Created: 0}
			m.ups[mid] = up
			if v, ok := Config.CoreConfig.GroupConfig[groupId]; ok {
				if v.BiliUps == nil {
					v.BiliUps = map[int64]Config.Up{}
				}
				v.BiliUps[mid] = Config.Up{Name: u.Data.Card.Name, Created: 0}
				Config.CoreConfig.GroupConfig[groupId] = v
			} else {
				if v.BiliUps == nil {
					v.BiliUps = map[int64]Config.Up{}
				}
				v = Config.CoreConfig.DefaultGroupConfig
				v.BiliUps[mid] = Config.Up{Name: u.Data.Card.Name, Created: 0}
				Config.CoreConfig.GroupConfig[groupId] = v
			}
			Config.Save()
		}

	}
	return
}

func (m *Manager) SearchUp(keyword string) (mid int64, err error) {
	res, err := m.r.Get("https://api.bilibili.com/x/web-interface/search/type?context=&search_type=bili_user&page=1&order=&category_id=&user_type=&order_sort=&changing=mid&__refresh__=true&_extra=&highlight=1&single_column=0&keyword=" + keyword)
	if err != nil {
		return
	}
	var result SearchResult
	err = res.Json(&result)
	if err != nil {
		return
	}
	mid = 0
	for _, v := range result.Data.Result {
		if v.IsUpuser == 1 {
			mid = v.Mid
			break
		}
	}
	if mid == 0 {
		err = errors.New("没有找到UP")
	}
	return
}
func (m *Manager) GetGroupsByMid(mid int64) (upName string, g []int64) {
	if v, ok := m.ups[mid]; ok {
		g = v.Groups
		upName = v.Name
	}
	return
}
