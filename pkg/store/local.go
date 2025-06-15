package store

import "sync"

type Local struct {
	m sync.Map
}

func (l *Local) Get(key string) (State, bool) {
	si, ok := l.m.Load(key)
	if !ok {
		return State{}, false
	}
	s, _ := si.(State)
	return s, true
}

func (l *Local) Set(key string, state State) {
	l.m.Store(key, state)
}

func NewLocalStore() *Local {
	return &Local{}
}
