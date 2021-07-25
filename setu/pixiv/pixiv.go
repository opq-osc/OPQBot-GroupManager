package pixiv

import (
	"OPQBot-QQGroupManager/Config"
	"OPQBot-QQGroupManager/Core"
	"OPQBot-QQGroupManager/setu/setucore"
	"OPQBot-QQGroupManager/utils"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/requests"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
)

var (
	log           *logrus.Entry
	codeVerifier  string
	CLIENT_ID     = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
	CLIENT_SECRET = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
	USER_AGENT    = "PixivAndroidApp/5.0.234 (Android 11; Pixel 5)"
	REDIRECT_URI  = "https://app-api.pixiv.net/web/v1/users/auth/pixiv/callback"
)

type Provider struct {
	c  Client
	db *gorm.DB
}

func (p *Provider) autoGetPic() {
	log.Info("自动获取Pixiv图片排行榜")
	result, err := p.c.GetDailyIllust()
	if err != nil {
		log.Warn(err)
	}
	addPicNum := 0
	for _, v := range result.Illusts {
		originPicUrl := ""
		if v.PageCount == 1 {
			originPicUrl = v.MetaSinglePage.OriginalImageURL
		} else {
			originPicUrl = v.MetaPages[0].ImageUrls.Original
		}
		var tag []string
		tagHasR18 := false
		for _, v1 := range v.Tags {
			tmp := ""
			if v1.TranslatedName != "" {
				tmp = v1.TranslatedName
			} else {
				tmp = v1.Name
			}
			if tmp == "R18" || tmp == "R-18" {
				tagHasR18 = true
			}
			tag = append(tag, tmp)
		}
		tmp := setucore.Pic{
			Id:             v.ID,
			Title:          v.Title,
			Author:         v.User.Name,
			AuthorID:       v.User.ID,
			OriginalPicUrl: originPicUrl,
			Tag:            strings.Join(tag, ","),
			R18:            v.XRestrict >= 1 || tagHasR18,
		}
		err := p.AddPicToDB(tmp)
		if err != nil {
			//log.Warn(err)
			continue
		}
		addPicNum += 1
	}
	log.Info("联网添加到本地数据库数据图片数量为:", addPicNum)
}
func (p *Provider) InitProvider(l *logrus.Entry, b *Core.Bot, db *gorm.DB) {
	log = l
	db.AutoMigrate(&setucore.Pic{})
	Config.Lock.RLock()
	debug := Config.CoreConfig.Debug
	p.c.Proxy = Config.CoreConfig.SetuConfig.PixivProxy
	p.c.refreshToken = Config.CoreConfig.SetuConfig.PixivRefreshToken
	autoGetPic := Config.CoreConfig.SetuConfig.AutoGetPic
	Config.Lock.RUnlock()
	p.db = db
	if debug {
		p.db = p.db.Debug()
	}
	if p.c.refreshToken == "" {
		p.c.GenerateLoginUrl()
	} else {
		err := p.c.RefreshToken()
		if err != nil {
			log.Error(err)
		}
	}
	if autoGetPic {
		err := b.BotCronManager.AddJob(-1, "setuAuto", "0 8 * * *", func() {
			p.autoGetPic()
		})
		if err != nil {
			log.Error(err)
		}
	}
	err := b.AddEvent(OPQBot.EventNameOnFriendMessage, func(qq int64, packet *OPQBot.FriendMsgPack) {
		if packet.FromUin != b.QQ {
			if strings.HasPrefix(packet.Content, "code=") && !p.c.Login {
				code := strings.TrimPrefix(packet.Content, "code=")
				err := p.c.LoginPixiv(code)
				if err != nil {
					log.Error(err)
				}
			}
			if p.c.Login && packet.Content == "px用户信息" {
				r, _ := requests.Get(p.c.loginInfo.Response.User.ProfileImageUrls.Px50X50)
				b.SendFriendPicMsg(packet.FromUin, fmt.Sprintf("用户名: %s (%s)", p.c.loginInfo.User.Name, p.c.loginInfo.User.ID), r.Content())
			}
		}
	})
	if err != nil {
		log.Error(err)
	}
}
func (p *Provider) SearchPicFromDB(word string, r18 bool, num int) (pics []setucore.Pic, e error) {
	if word == "" {
		e = p.db.Where("r18 = ? AND last_send_time < ?", "%"+word+"%", r18, time.Now().Unix()-1800).Limit(num).Order("last_send_time asc").Find(&pics).Error
		return
	}
	e = p.db.Where("tag LIKE ? AND r18 = ? AND last_send_time < ?", "%"+word+"%", r18, time.Now().Unix()-1800).Limit(num).Order("last_send_time asc").Find(&pics).Error
	return
}
func (p *Provider) AddPicToDB(pic setucore.Pic) error {
	var num int64
	p.db.Model(&pic).Where("id = ?", pic.Id).Count(&num)
	if num > 0 {
		return errors.New("该图片在数据库中已存在！")
	}
	return p.db.Create(&pic).Error
}
func (p *Provider) SetPicSendTime(pics []setucore.Pic) {
	var sendPicId []int
	for _, v := range pics {
		sendPicId = append(sendPicId, v.Id)
	}
	p.db.Model(&setucore.Pic{}).Where("id IN ?", sendPicId).Updates(&setucore.Pic{LastSendTime: time.Now().Unix()})
}
func (p *Provider) SearchPic(word string, r18 bool, num int) ([]setucore.Pic, error) {
	dbPic, err := p.SearchPicFromDB(word, r18, num)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	if len(dbPic) < num {
		log.Info("本地数据库数据量不够，联网获取中...")
		result, err := p.c.SearchIllust(word)
		if err != nil {
			log.Warn(err)
		}
		addPicNum := 0
		for _, v := range result.Illusts {
			originPicUrl := ""
			if v.PageCount == 1 {
				originPicUrl = v.MetaSinglePage.OriginalImageURL
			} else {
				originPicUrl = v.MetaPages[0].ImageUrls.Original
			}
			var tag []string
			tagHasR18 := false
			for _, v1 := range v.Tags {
				tmp := ""
				if v1.TranslatedName != "" {
					tmp = v1.TranslatedName
				} else {
					tmp = v1.Name
				}
				if tmp == "R18" || tmp == "R-18" {
					tagHasR18 = true
				}
				tag = append(tag, tmp)
			}
			tmp := setucore.Pic{
				Id:             v.ID,
				Title:          v.Title,
				Author:         v.User.Name,
				AuthorID:       v.User.ID,
				OriginalPicUrl: originPicUrl,
				Tag:            strings.Join(tag, ","),
				R18:            v.XRestrict >= 1 || tagHasR18,
			}
			err := p.AddPicToDB(tmp)
			if err != nil {
				//log.Warn(err)
				continue
			}
			addPicNum += 1
		}
		log.Info("联网添加到本地数据库数据关于", word, "的记录数量为:", addPicNum)
		dbPic, err = p.SearchPicFromDB(word, r18, num)
		if err != nil {
			log.Warn(err)
			return nil, err
		}
	}
	if len(dbPic) > 0 {
		p.SetPicSendTime(dbPic)
	}
	return dbPic, nil
}

type Client struct {
	Proxy            string
	Login            bool
	tokenExpiresTime time.Time
	refreshToken     string
	Token            string
	loginInfo        LoginSuccessResult
}

func (c *Client) RefreshToken() error {
	// 尝试通过 refreshToken 获取Token
	if time.Now().After(c.tokenExpiresTime) {

		req := requests.Requests()
		if c.Proxy != "" {
			req.Proxy(c.Proxy)
		}
		r, err := req.Post("https://oauth.secure.pixiv.net/auth/token",
			requests.Datas{
				"client_id":      CLIENT_ID,
				"client_secret":  CLIENT_SECRET,
				"grant_type":     "refresh_token",
				"refresh_token":  c.refreshToken,
				"get_secure_url": "1",
			},
			requests.Header{
				"User-Agent": USER_AGENT,
				"host":       "oauth.secure.pixiv.net",
			})
		if err != nil {
			return err
		}
		tmp := map[string]interface{}{}
		err = r.Json(&tmp)
		if err != nil {
			return err
		}
		if _, ok := tmp["has_error"]; ok {
			result := ErrResult{}
			err = r.Json(&result)
			if err != nil {
				return err
			}
			return errors.New(fmt.Sprintf("[%d]%s", result.Errors.System.Code, result.Errors.System.Message))
		}
		result := LoginSuccessResult{}
		err = r.Json(&result)
		if err != nil {
			return err
		}
		c.loginInfo = result
		c.Login = true
		c.Token = result.AccessToken
		c.refreshToken = result.RefreshToken
		Config.Lock.Lock()
		Config.CoreConfig.SetuConfig.PixivRefreshToken = result.RefreshToken
		Config.Save()
		Config.Lock.Unlock()

		c.tokenExpiresTime = time.Now().Add(time.Second * time.Duration(result.ExpiresIn))
		log.Println("登录成功")
	}
	return nil
}
func (c *Client) GetUserInfo() LoginSuccessResult {
	return c.loginInfo
}
func (c *Client) LoginPixiv(code string) error {
	// 尝试通过 Code 获取 Refresh Token
	if c.refreshToken != "" {
		return c.RefreshToken()
	}
	req := requests.Requests()
	if c.Proxy != "" {
		req.Proxy(c.Proxy)
	}
	r, err := req.Post("https://oauth.secure.pixiv.net/auth/token",
		requests.Datas{
			"client_id":      CLIENT_ID,
			"client_secret":  CLIENT_SECRET,
			"code":           code,
			"code_verifier":  codeVerifier,
			"grant_type":     "authorization_code",
			"include_policy": "true",
			"redirect_uri":   REDIRECT_URI,
		},
		requests.Header{
			"User-Agent": USER_AGENT,
			"host":       "oauth.secure.pixiv.net",
		})
	if err != nil {
		return err
	}
	tmp := map[string]interface{}{}
	err = r.Json(&tmp)
	if err != nil {
		return err
	}
	if _, ok := tmp["has_error"]; ok {
		result := ErrResult{}
		err = r.Json(&result)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("[%d]%s", result.Errors.System.Code, result.Errors.System.Message))
	}
	result := LoginSuccessResult{}
	err = r.Json(&result)
	if err != nil {
		return err
	}
	c.loginInfo = result
	c.Login = true
	c.Token = result.AccessToken
	c.refreshToken = result.RefreshToken
	Config.Lock.Lock()
	Config.CoreConfig.SetuConfig.PixivRefreshToken = result.RefreshToken
	Config.Save()
	Config.Lock.Unlock()
	c.tokenExpiresTime = time.Now().Add(time.Second * time.Duration(result.ExpiresIn))
	log.Println("登录成功")
	return nil
}
func (c *Client) GenerateLoginUrl() {
	codeVerifier = Base64UrlSafeEncode([]byte(utils.RandomString(32)))
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	sum := h.Sum(nil)
	codeChallenge := Base64UrlSafeEncode(sum)
	log.Println("未登录! 请登录以下网址，然后将链接中的code参数私信发送给机器人完成登录。请发送\"code=xxxxxxx\"")
	log.Println("https://app-api.pixiv.net/web/v1/login?code_challenge=" + codeChallenge + "&code_challenge_method=S256&client=pixiv-android")
}
func (c *Client) GetHeader() requests.Header {
	err := c.RefreshToken()
	if err != nil {
		log.Warn(err)
	}
	if c.Token != "" {
		return requests.Header{
			"host":            "app-api.pixiv.net",
			"app-os":          "ios",
			"Accept-Language": "zh-cn",
			"User-Agent":      "PixivIOSApp/7.13.3 (iOS 14.6; iPhone13,2)",
			"Authorization":   fmt.Sprintf("Bearer %s", c.Token),
		}
	} else {
		return requests.Header{
			"Accept-Language": "zh-cn",
			"User-Agent":      USER_AGENT,
		}
	}

}
func (c *Client) SearchIllust(word string) (result IllustResult, err error) {
	var res *requests.Response
	req := requests.Requests()
	if c.Proxy != "" {
		req.Proxy(c.Proxy)
	}
	res, err = req.Get(fmt.Sprintf("https://app-api.pixiv.net/v1/search/popular-preview/illust?word=%s&search_target=partial_match_for_tags&sort=date_desc&filter=for_ios", word), c.GetHeader())
	if err != nil {
		return
	}
	err = res.Json(&result)
	return
}
func (c *Client) GetDailyIllust() (result IllustResult, err error) {
	var res *requests.Response
	req := requests.Requests()
	if c.Proxy != "" {
		req.Proxy(c.Proxy)
	}
	res, err = req.Get("https://app-api.pixiv.net/v1/illust/ranking?mode=day&filter=for_ios", c.GetHeader())
	if err != nil {
		return
	}
	err = res.Json(&result)
	return
}
func Base64URLDecode(data string) ([]byte, error) {
	var missing = (4 - len(data)%4) % 4
	data += strings.Repeat("=", missing)
	res, err := base64.URLEncoding.DecodeString(data)
	fmt.Println("  decodebase64urlsafe is :", string(res), err)
	return base64.URLEncoding.DecodeString(data)
}
func Base64UrlSafeEncode(source []byte) string {
	// Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-') and trailing '=' are removed.
	bytearr := base64.StdEncoding.EncodeToString(source)
	safeurl := strings.Replace(string(bytearr), "/", "_", -1)
	safeurl = strings.Replace(safeurl, "+", "-", -1)
	safeurl = strings.Replace(safeurl, "=", "", -1)
	return safeurl
}
