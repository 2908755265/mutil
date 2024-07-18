package priority

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type PriorityFunc func(score int64) int64

type Task struct {
	score     int64
	realScore int64
	Data      any
	Run       func() error
	ctx       context.Context
	cancel    func()
}

func DefaultPriorityFunc(score int64) int64 {
	return score - time.Now().Unix()
}

func NewTask(score int64, pf PriorityFunc, data any, run func() error, cancelCtx context.Context, cancel func()) *Task {
	return &Task{
		score:     score,
		realScore: pf(score),
		Data:      data,
		Run:       run,
		ctx:       cancelCtx,
		cancel:    cancel,
	}
}

func NewDefaultTask(score int64, data any, run func() error) *Task {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return NewTask(score, DefaultPriorityFunc, data, run, ctx, cancelFunc)
}

type PriorityLock struct {
	listLock     *sync.Mutex
	tasks        *list.List
	dispatchLock *sync.Mutex
	hasNewTask   chan struct{}
}

func NewPriorityLock(size uint) *PriorityLock {
	return &PriorityLock{
		listLock:     &sync.Mutex{},
		tasks:        list.New(),
		dispatchLock: &sync.Mutex{},
		hasNewTask:   make(chan struct{}, size),
	}
}

func (l *PriorityLock) addTask(tsk *Task) {
	l.listLock.Lock()
	defer l.listLock.Unlock()
	back := l.tasks.Back()
	if back == nil {
		l.tasks.PushBack(tsk)
	} else {
		for {
			task := back.Value.(*Task)
			if task.realScore >= tsk.realScore {
				l.tasks.InsertAfter(tsk, back)
				break
			} else {
				next := back.Prev()
				if next == nil {
					l.tasks.PushFront(tsk)
					break
				} else {
					back = next
				}
			}
		}
	}
	l.hasNewTask <- struct{}{}
}

func (l *PriorityLock) pullTask() *Task {
	l.listLock.Lock()
	defer l.listLock.Unlock()
	if l.tasks.Len() > 0 {
		front := l.tasks.Front()
		l.tasks.Remove(front)
		return front.Value.(*Task)
	}
	return nil
}

func (l *PriorityLock) dispatch() {
	for {
		select {
		case <-l.hasNewTask:
			l.dispatchLock.Lock()
			task := l.pullTask()
			if task == nil {
				l.dispatchLock.Unlock()
				return
			}
			task.cancel()
		}
	}
}

func (l *PriorityLock) Start() {
	go l.dispatch()
}

func (l *PriorityLock) Stop() {
	close(l.hasNewTask)
}

func (l *PriorityLock) Lock(tsk *Task) {
	l.addTask(tsk)
	select {
	case <-tsk.ctx.Done():
	}
}

func (l *PriorityLock) UnLock() {
	l.dispatchLock.Unlock()
}
