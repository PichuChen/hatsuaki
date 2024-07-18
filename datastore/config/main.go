package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	// Domain 是這個服務的域名，舉例來說 @alice@example.com 的 domain 就是 example.com
	Domain string `json:"domain"`
}

var runningConfig Config

func GetDomain() string {
	return runningConfig.Domain
}

func LoadConfig(filepath string) error {
	f, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(f, &runningConfig)
	if err != nil {
		return err
	}
	return nil
}

func SaveConfig(filepath string) error {
	f, err := json.MarshalIndent(runningConfig, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, f, 0644)
	if err != nil {
		return err
	}
	return nil
}
