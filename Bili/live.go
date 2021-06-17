package bili

import (
	"errors"
	"fmt"
	"github.com/gaoyanpao/biliLiveHelper"
	"sync"
)

type LiveManager struct {
	biliClient map[int]*Client
	lock       sync.RWMutex
}

type Client struct {
	C         *biliLiveHelper.Client
	OnDanmu   func(ctx *biliLiveHelper.Context)
	OnGift    func(ctx *biliLiveHelper.Context)
	OnWelcome func(ctx *biliLiveHelper.Context)
}

func NewLiveManager() *LiveManager {
	return &LiveManager{biliClient: make(map[int]*Client), lock: sync.RWMutex{}}
}

func (l *LiveManager) AddClient(RoomId int) (c *Client, e error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	client := biliLiveHelper.NewClient(RoomId)
	if client == nil {
		e = errors.New("Client Err ")
		return
	}
	client.RegHandleFunc(biliLiveHelper.CmdDanmuMsg, func(ctx *biliLiveHelper.Context) bool {
		if c.OnDanmu != nil {
			c.OnDanmu(ctx)
		}
		return false
	})
	client.RegHandleFunc(biliLiveHelper.CmdSendGift, func(ctx *biliLiveHelper.Context) bool {
		if c.OnGift != nil {
			c.OnGift(ctx)
		}
		return false
	})
	client.RegHandleFunc(biliLiveHelper.CmdWelcome, func(ctx *biliLiveHelper.Context) bool {
		if c.OnWelcome != nil {
			c.OnWelcome(ctx)
		}
		return false
	})
	c = &Client{C: client}
	if _, ok := l.biliClient[RoomId]; ok {
		e = errors.New("已加入了Room")
		return
	}
	l.biliClient[RoomId] = c
	return
}

func (l *LiveManager) RemoveClient(RoomId int) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if v, ok := l.biliClient[RoomId]; ok {
		if v.C.IsConnected() {
			v.ExitRoom()
		}
		delete(l.biliClient, RoomId)
		return nil
	} else {
		return errors.New("Client不存在! ")
	}
}

func (c *Client) RegisterDanmuFunc(f func(ctx *biliLiveHelper.Context)) {
	c.OnDanmu = f
}

func (c *Client) RegisterGiftFunc(f func(ctx *biliLiveHelper.Context)) {
	c.OnGift = f
}
func (c *Client) RegisterWelcomeFunc(f func(ctx *biliLiveHelper.Context)) {
	c.OnWelcome = f
}

func (c *Client) GetRoomInfo() biliLiveHelper.SimpleRoomInfo {
	return c.C.SimpleRoomInfo
}

func (c *Client) ExitRoom() error {
	return c.C.Conn.Close()
}

// Start 阻塞
func (c *Client) Start() error {

	return c.C.StartListen()
}
func GetLiveStatusString(status int) string {
	switch status {
	case 0:
		return "未开播"
	case 1:
		return "直播中"
	case 2:
		return "轮播中"
	default:
		return fmt.Sprintf("未知状态%d", status)
	}
}
