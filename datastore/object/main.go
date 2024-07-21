package object

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"sync"
	"time"
)

type Object map[string]interface{}

var datastore = &sync.Map{}

func LoadObject(filepath string) error {
	slog.Debug("object.Load", "info", "load objects")

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
		a := Object(m)
		tmpDatastore.Store(k, &a)
	}

	// old datastore should be garbage collected
	datastore = &tmpDatastore
	slog.Info("object.Load", "info", "objects loaded")
	return nil
}

func SaveObject(filepath string) error {
	slog.Debug("object.Save", "info", "save objects", "filepath", filepath)

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

	slog.Info("object.Save", "info", "objects saved")
	return nil
}

func FindObjectByID(id string) (*Object, error) {
	if v, ok := datastore.Load(id); ok {
		return v.(*Object), nil
	}
	return nil, fmt.Errorf("object not found")
}

func GenerateUUIDv7() string {
	// UUIDv7
	var buf [16]byte
	rand.Read(buf[:])
	t := big.NewInt(time.Now().UnixMilli())
	t.FillBytes(buf[:6])
	buf[6] = 0x70 | (buf[6] & 0x0f)
	buf[8] = 0x80 | (buf[8] & 0x3f)
	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}
