package actor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/pichuchen/hatsuaki/activitypub/signature"
	"golang.org/x/crypto/bcrypt"
)

type Actor map[string]interface{}

var datastore = &sync.Map{}

func LoadActor(filepath string) error {
	slog.Debug("actor.Load", "info", "load actors")

	f, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	tmpMap := map[string]interface{}{}
	tmpDatastore := sync.Map{}

	err = json.Unmarshal(f, &tmpMap)
	if err != nil {
		return err
	}

	for k, v := range tmpMap {
		m := v.(map[string]interface{})
		a := Actor(m)
		tmpDatastore.Store(k, &a)
	}

	if _, ok := tmpDatastore.Load("instance.actor"); !ok {
		// 如果讀取了檔案，但是裡面卻沒有 instance.actor 的話 (可能被刪掉了)
		// initial instance.actor
		InitActorDatastore()
	}
	// old datastore should be garbage collected
	datastore = &tmpDatastore
	slog.Info("actor.Load", "info", "actors loaded")
	return nil
}

func SaveActor(filepath string) error {
	slog.Debug("actor.Save", "info", "save actors", "filepath", filepath)

	tmpMap := map[string]interface{}{}
	datastore.Range(func(k, v interface{}) bool {
		tmpMap[k.(string)] = v
		return true
	})

	f, err := json.MarshalIndent(tmpMap, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, f, 0644)
	if err != nil {
		return err
	}

	slog.Info("actor.Save", "info", "actors saved")
	return nil
}

func InitActorDatastore() {
	if datastore == nil {
		datastore = &sync.Map{}
	}
	datastore.Store("instance.actor", &Actor{
		"username":   "instance.actor",
		"privateKey": signature.GeneratePrivateKey(),
	})
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

// 會以 PEM 格式回傳 RSA Private Key
func (a *Actor) GetPrivateKey() string {
	key, ok := (*a)["privateKey"].(string)
	if !ok {
		slog.Warn("actor.GetPrivateKey", "error", "privateKey not found")
		// 我們在這邊產生一個新的
		key = signature.GeneratePrivateKey()
		(*a)["privateKey"] = key
	}
	return key
}

// 會以 PEM 格式回傳 RSA Public Key
func (a *Actor) GetPublicKey() string {
	p := a.GetPrivateKey()
	return string(signature.Pubout([]byte(p)))
}

func NewActor(username string) *Actor {
	a := &Actor{
		"username":   username,
		"privateKey": signature.GeneratePrivateKey(),
	}
	datastore.Store(username, a)
	return a
}

func UpdatePassword(username, password string) error {
	a, err := FindActorByUsername(username)
	if err != nil {
		slog.Error("actor.UpdatePassword", "error", err)
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("actor.UpdatePassword", "error", err)
		return err
	}
	(*a)["hashedPassword"] = string(hashedPassword)
	return nil
}
