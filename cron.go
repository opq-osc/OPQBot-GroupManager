package main

import (
	"errors"
	"strconv"
	"sync"

	"github.com/robfig/cron/v3"
)

type BotCron struct {
	c      *cron.Cron
	jobs   map[string]cron.EntryID
	locker sync.Mutex
}

func NewBotCronManager() BotCron {
	return BotCron{c: cron.New(), jobs: map[string]cron.EntryID{}, locker: sync.Mutex{}}
}

func (m *BotCron) AddJob(qqGroup int64, jobName string, spec string, cmd func()) error {
	Name := strconv.FormatInt(qqGroup, 10) + "-" + jobName
	m.locker.Lock()
	defer m.locker.Unlock()
	if _, ok := m.jobs[Name]; ok {
		return errors.New("Job已经存在了，不能重复添加!")
	} else {
		if id, err := m.c.AddFunc(spec, cmd); err != nil {
			return err
		} else {
			m.jobs[Name] = id
			return nil
		}
	}
}

func (m *BotCron) Remove(qqGroup int64, jobName string) error {
	Name := strconv.FormatInt(qqGroup, 10) + "-" + jobName
	m.locker.Lock()
	defer m.locker.Unlock()
	if v, ok := m.jobs[Name]; ok {
		m.c.Remove(v)
		delete(m.jobs, Name)
		return nil
	} else {
		return errors.New("Job不存在!")
	}
}

func (m *BotCron) List() map[string]cron.EntryID {
	m.locker.Lock()
	defer m.locker.Unlock()
	return m.jobs
}

func (m *BotCron) Start() {
	m.c.Start()
}
