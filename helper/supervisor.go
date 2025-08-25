package helper

import (
	"sync"
)

type Supervisor struct {
	tasks map[string]*Task
}

var (
	s *Supervisor
)

func init() {
	if s == nil {
		s = &Supervisor{
			tasks: make(map[string]*Task),
		}
	}
}

func Push(name string, task func()) {
	s.Push(name, task)
}

func (s *Supervisor) Push(name string, task func()) {
	t := Task{
		handle:    task,
		retry:     0,
		max_retry: 3,
		running:   false,
	}

	s.tasks[name] = &t

	go t.Run()
}

type Task struct {
	handle    func()
	retry     int // 当前第几次重试
	max_retry int // 最大重试次数
	mu        sync.Mutex
	running   bool
}

func (t *Task) Run() {
	defer func() {
		if e := recover(); e != nil {
			// 异常退出重试
			t.running = false
			t.retry += 1
			t.restart()
		}
	}()

	// 防止出现并发重试
	if !t.running {
		t.mu.Lock()
		defer t.mu.Unlock()

		if !t.running {
			t.running = true
			t.handle()
		}
	}
}

func (t *Task) restart() {
	// 停止重试
	if t.retry > t.max_retry {
		return
	}

	go t.Run()
}
