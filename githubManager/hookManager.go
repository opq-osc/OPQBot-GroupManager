package githubManager

import (
	"OPQBot-QQGroupManager/Config"
	"errors"
	"fmt"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/kataras/iris/v12"
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/requests"
	"log"
	"strings"
	"sync"
)

type Manager struct {
	github     map[string]Config.Repo
	githubLock sync.RWMutex
	b          *OPQBot.BotManager
}

func (m *Manager) DelRepo(repo string, groupId int64) error {
	m.githubLock.Lock()
	defer m.githubLock.Unlock()
	if v, ok := m.github[repo]; ok {
		for i, v1 := range v.Groups {
			if v1 == groupId {
				v.Groups = append(v.Groups[:i], v.Groups[i+1:]...)
				m.github[repo] = v
				if len(v.Groups) == 0 {
					delete(m.github, repo)
				}
				break
			}
		}
		m.Save()
		return nil

	} else {
		return errors.New("订阅不存在!无法删除 ")
	}
}
func (m *Manager) AddRepo(repo string, Secret string, groupId int64) error {
	m.githubLock.Lock()
	defer m.githubLock.Unlock()

	if v, ok := m.github[repo]; ok {
		if v.Secret != Secret {
			v.Secret = Secret
			hook, _ := github.New(github.Options.Secret(Secret))
			v.WebHook = hook
		}
		for _, v1 := range v.Groups {
			if v1 == groupId {
				m.Save()
				return errors.New("已经订阅过了")
			}
		}
		v.Groups = append(v.Groups, groupId)
		m.github[repo] = v
	} else {
		if Secret == "" {
			return errors.New("需要Secret, 请直接私聊我发送Secret")
		}
		hook, _ := github.New(github.Options.Secret(Secret))
		m.github[repo] = Config.Repo{
			WebHook: hook,
			Groups:  []int64{groupId},
			Secret:  Secret,
		}
	}
	err := m.Save()
	if err != nil {
		return err
	}
	return nil
}

func NewManager(app *iris.Application, bot *OPQBot.BotManager) Manager {
	Config.Lock.RLock()
	defer Config.Lock.RUnlock()
	g := map[string]Config.Repo{}
	for k, v := range Config.CoreConfig.GithubSub {
		hook, _ := github.New(github.Options.Secret(v.Secret))
		v.WebHook = hook
		g[k] = v
	}
	m := Manager{github: g, githubLock: sync.RWMutex{}, b: bot}
	app.Any("github/webhook/{root:path}", func(ctx iris.Context) {
		h, err := m.GetRepo(ctx.Params().GetString("root"))
		if err != nil {
			ctx.StatusCode(404)
			return
		}
		payload, err := h.WebHook.Parse((*ctx).Request(), github.RepositoryEvent, github.PushEvent, github.PingEvent, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			log.Println(err)
			if err == github.ErrEventNotFound {
				ctx.StatusCode(404)
				return
			}
			if err == github.ErrHMACVerificationFailed {
				ctx.StatusCode(502)
				return
			}
		}
		switch v := payload.(type) {
		case github.PingPayload:
			log.Println(v)
		case github.RepositoryPayload:
			switch v.Action {
			case "created":
				r, _ := requests.Get(v.Sender.AvatarURL)
				for _, v1 := range h.Groups {
					m.b.SendGroupPicMsg(v1, fmt.Sprintf("%s在%s发布了新的仓库: %s\n欢迎Star哟", v.Sender.Login, v.Organization.Login, v.Repository.FullName), r.Content())
				}
			}

		case github.PushPayload:
			var commitString []string
			for _, v1 := range v.Commits {
				commitString = append(commitString, fmt.Sprintf("[%s] %s", v1.Timestamp, v1.Message))
			}
			if len(commitString) == 0 {
				return
			}
			r, _ := requests.Get(v.Sender.AvatarURL)
			for _, v1 := range h.Groups {
				m.b.SendGroupPicMsg(v1, fmt.Sprintf("%s\n%s发起了Push\nCommit:\n%s", v.Repository.FullName, v.Pusher.Name, strings.Join(commitString, "\n")), r.Content())
			}
		case github.ReleasePayload:
			r, _ := requests.Get(v.Sender.AvatarURL)
			switch v.Action {
			case "published":
				for _, v1 := range h.Groups {
					m.b.SendGroupPicMsg(v1, fmt.Sprintf("%s\n%s发布了新版本:\n%s", v.Repository.FullName, v.Sender.Login, v.Release.TagName), r.Content())
				}
			default:

			}

		case github.PullRequestPayload:
			r, _ := requests.Get(v.PullRequest.User.AvatarURL)
			msg := ""
			switch v.Action {
			case "closed":
				msg = fmt.Sprintf("%s\n%s关闭了PR:%s to %s", v.Repository.FullName, v.PullRequest.User.Login, v.PullRequest.Head.Label, v.PullRequest.Base.Label)
			case "opened":
				msg = fmt.Sprintf("%s\n%s打开了PR:%s to %s", v.Repository.FullName, v.PullRequest.User.Login, v.PullRequest.Head.Label, v.PullRequest.Base.Label)
			default:
				ctx.StatusCode(503)
				return
			}
			for _, v1 := range h.Groups {
				m.b.SendGroupPicMsg(v1, msg, r.Content())
			}

		}
	})
	return m
}
func (m *Manager) GetGroupSubList(groupId int64) (r map[string]Config.Repo) {
	m.githubLock.RLock()
	defer m.githubLock.RUnlock()
	r = map[string]Config.Repo{}
	for k, v := range m.github {
		for _, v1 := range v.Groups {
			if v1 == groupId {
				r[k] = v
			}
		}
	}
	return
}

func (m *Manager) GetRepo(repo string) (r Config.Repo, err error) {
	m.githubLock.RLock()
	defer m.githubLock.RUnlock()
	var ok bool
	if r, ok = m.github[repo]; ok {
		return
	} else {
		err = errors.New("没有订阅该Repo")
		return
	}

}

func (m *Manager) Save() error {
	Config.Lock.Lock()
	defer Config.Lock.Unlock()
	Config.CoreConfig.GithubSub = m.github
	return Config.Save()
}
