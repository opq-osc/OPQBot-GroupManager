package Chat

import (
	"OPQBot-QQGroupManager/Config"
	"errors"

	"github.com/sirupsen/logrus"
)

type Manager struct {
	SelectCore string
}

var (
	log       *logrus.Entry
	Providers map[string]ChatCore
)

type ChatCore interface {
	Init(l *logrus.Entry) error
	GetAnswer(question string, GroupId, userId int64) (string, []byte)
	AddAnswer(question, answer string, GroupId, userId int64) error
	SetReplace(regexp string, target string) error
}

func StartChatCore(l *logrus.Entry) Manager {
	log = l
	Config.Lock.RLock()
	tmp := Config.CoreConfig.SelectChatCore
	Config.Lock.RUnlock()
	for k, v := range Providers {
		tmp := log.WithField("Provider", k)
		tmp.Info("载入中")
		err := v.Init(tmp)
		if err != nil {
			tmp.Error(err)
			delete(Providers, k)
		}
		tmp.Info("载入成功")
	}
	if _, ok := Providers[tmp]; ok {
		return Manager{SelectCore: tmp}
	} else {
		return Manager{SelectCore: ""}
	}
}

func init() {
	Providers = make(map[string]ChatCore)
}
func (m *Manager) Learn(Question, Answer string, GroupId, From int64) error {
	if v, ok := Providers["local"]; ok {
		return v.AddAnswer(Question, Answer, GroupId, From)
	} else {
		return errors.New("本地聊天系统出现故障")
	}
}
func (m *Manager) GetChatDB() string {
	return m.SelectCore
}
func (m *Manager) SetChatDB(db string) error {
	if _, ok := Providers[db]; ok {
		m.SelectCore = db
		Config.Lock.Lock()
		Config.CoreConfig.SelectChatCore = db
		Config.Save()
		Config.Lock.Unlock()
		return nil
	}
	return errors.New("数据库不存在")
}
func (m *Manager) GetAnswer(question string, groupId, userId int64) (string, []byte, error) {
	// 查找本地对话数据库
	if v, ok := Providers["local"]; ok {
		if answer, pic := v.GetAnswer(question, groupId, userId); answer != "" || pic != nil {
			return answer, pic, nil
		}
	}
	// 联网查询默认对话数据库
	if v, ok := Providers[m.SelectCore]; ok {
		if answer, pic := v.GetAnswer(question, groupId, userId); answer != "" || pic != nil {
			return answer, pic, nil
		}
	}
	// 遍历其他数据库
	for k, v := range Providers {
		if k == "local" || k == m.SelectCore {
			continue
		}
		if answer, pic := v.GetAnswer(question, groupId, userId); answer != "" {
			return answer, pic, nil
		}
	}
	return "", nil, errors.New("没有找到对话记录，无法回答")
}
func Register(name string, core ChatCore) error {
	if _, ok := Providers[name]; ok {
		return errors.New("Core已经注册了 ")
	} else {
		Providers[name] = core
	}
	return nil
}
