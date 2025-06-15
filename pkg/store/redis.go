package store

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.ClusterClient
}

func (l *Redis) Get(key string) (State, bool) {
	//TODO: Pass context properly
	data, err := l.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return State{}, false
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, false
	}
	return s, true
}

func (l *Redis) Set(key string, state State) {
	data, err := json.Marshal(state)
	if err != nil {
		return
	}
	//TODO: Pass context properly
	//TODO: Set expiration time properly
	l.client.Set(context.Background(), key, data, 0).Err()
}

func NewRedisStore(client *redis.ClusterClient) *Redis {
	return &Redis{
		client: client,
	}
}
