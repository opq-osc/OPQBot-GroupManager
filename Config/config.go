package Config

import (
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
	DefaultGroupConfig GroupConfig
	SuperAdminUin      int64
	WhiteGroupList     []int64
	BlackGroupList     []int64
	GroupConfig        map[int64]GroupConfig
	UserData           map[int64]UserData
}
type UserData struct {
	LastSignDay int
	LastZanDay  int
	Count       int
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
	Welcome            string
	JoinVerifyType     int
	Job                map[string]Job
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
