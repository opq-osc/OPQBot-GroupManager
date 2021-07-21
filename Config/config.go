package Config

import (
	"github.com/go-playground/webhooks/v6/github"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type CoreConfigStruct struct {
	OPQWebConfig struct {
		CSRF     string
		Host     string
		Port     int
		Username string
		Password string
		Enable   bool
	}
	OPQBotConfig struct {
		Url string
		QQ  int64
	}
	BiliLive           bool
	YiQing             bool
	ReverseProxy       string
	DefaultGroupConfig GroupConfig
	SuperAdminUin      int64
	WhiteGroupList     []int64
	BlackGroupList     []int64
	GroupConfig        map[int64]GroupConfig
	UserData           map[int64]UserData
	GithubSub          map[string]Repo
	LogLevel           string
}
type Repo struct {
	Secret  string
	WebHook *github.Webhook
	Groups  []int64
}
type UserData struct {
	LastSignDay int
	LastZanDay  int
	Count       int
	SteamShare  string
}
type Job struct {
	Cron    string
	Type    int
	Title   string
	Content string
}
type GroupConfig struct {
	Enable             bool
	AdminUin           int64
	Menu               string
	MenuKeyWord        string
	ShutUpWord         string
	ShutUpTime         int
	JoinVerifyTime     int
	JoinAutoShutUpTime int
	Zan                bool
	SignIn             bool
	Bili               bool
	BiliUps            map[int64]Up
	Fanjus             map[int64]Fanju
	Welcome            string
	JoinVerifyType     int
	Job                map[string]Job
}
type Up struct {
	Name    string
	Created int64
	UserId  int64
}
type Fanju struct {
	Title  string
	Id     int64
	UserId int64
}

var (
	CoreConfig = &CoreConfigStruct{}
	Lock       = sync.RWMutex{}
)

func Save() error {
	b, err := yaml.Marshal(CoreConfig)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./config.yaml", b, 0777)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) == 2 && os.Args[1] == "first" {
		b, err := yaml.Marshal(CoreConfig)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("./config.yaml.example", b, 0777)
		if err != nil {
			panic(err)
		}
		panic("已将默认配置文件写出")
	}
	b, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Println("读取配置文件失败")
		panic(err)
	}
	err = yaml.Unmarshal(b, &CoreConfig)
	if err != nil {
		log.Println("读取配置文件失败")
		panic(err)
	}
	Save()
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					b, err := ioutil.ReadFile("./config.yaml")
					if err != nil {
						log.Println("读取配置文件失败")
						break
					}
					Lock.Lock()
					err = yaml.Unmarshal(b, &CoreConfig)
					Lock.Unlock()
					if err != nil {
						log.Println("读取配置文件失败")
						break
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
}
