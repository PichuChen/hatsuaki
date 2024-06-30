package actor

import (
	"fmt"
	"log/slog"
	"sync"
)

type Actor map[string]interface{}

var datastore = sync.Map{}

func Load() {
	slog.Debug("actor.Load", "info", "load actors")
	datastore.Store("alice", &Actor{"username": "alice"})
	datastore.Store("bob", &Actor{"username": "bob"})

	slog.Info("actor.Load", "info", "actors loaded")
}

func FindActorByUsername(username string) (actor *Actor, err error) {
	slog.Info("actor.FindActorByUsername", "username", username)
	if a, ok := datastore.Load(username); ok {
		actor = a.(*Actor)
		return actor, nil
	}
	return nil, fmt.Errorf("actor not found")
}

func (a *Actor) GetUsername() string {
	return (*a)["username"].(string)
}
