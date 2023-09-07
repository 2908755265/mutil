package msync

import "sync"

type MapInterface[KEY comparable, VAL any] interface {
	Load(KEY) (VAL, bool)
	Store(KEY, VAL)
	LoadOrStore(KEY, VAL) (actual VAL, loaded bool)
	LoadAndDelete(KEY) (v VAL, loaded bool)
	Delete(KEY)
	Range(func(KEY, VAL) (shouldContinue bool))
}

type syncMap[KEY comparable, VAL any] struct {
	zero VAL
	sm   *sync.Map
}

func (m *syncMap[KEY, VAL]) Load(k KEY) (VAL, bool) {
	v, ok := m.sm.Load(k)
	if ok {
		return v.(VAL), ok
	}
	return m.zero, ok
}

func (m *syncMap[KEY, VAL]) Store(k KEY, v VAL) {
	m.sm.Store(k, v)
}

func (m *syncMap[KEY, VAL]) LoadOrStore(k KEY, v VAL) (actual VAL, loaded bool) {
	store, ok := m.sm.LoadOrStore(k, v)
	return store.(VAL), ok
}

func (m *syncMap[KEY, VAL]) LoadAndDelete(k KEY) (v VAL, loaded bool) {
	value, ok := m.sm.LoadAndDelete(k)
	if !ok {
		return m.zero, ok
	}
	return value.(VAL), ok
}

func (m *syncMap[KEY, VAL]) Delete(k KEY) {
	m.sm.Delete(k)
}

func (m *syncMap[KEY, VAL]) Range(fn func(KEY, VAL) (shouldContinue bool)) {
	m.sm.Range(func(key, value any) bool {
		return fn(key.(KEY), value.(VAL))
	})
}

func NewSyncMap[KEY comparable, VAL any]() MapInterface[KEY, VAL] {
	var v VAL
	return &syncMap[KEY, VAL]{
		sm:   &sync.Map{},
		zero: v,
	}
}
