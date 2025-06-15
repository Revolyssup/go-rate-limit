package store

type State struct {
	LastReq int64
	Excess  int64
}

type Store interface {
	Get(key string) (State, bool)
	Set(key string, state State)
}
