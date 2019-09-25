package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DebugLevel int    `json:"debug"`
	Universe   int    `json:"universe"`
	Address    int    `json:"address"`
	Interface  string `json:"interface"`
	Protocol   string `json:"protocol"`
}

func Load() (cfg Config, err error) {
	f, err := os.Open("config.json")
	if err != nil {
		return
	}
	defer f.Close()

	j := json.NewDecoder(f)
	err = j.Decode(&cfg)
	return
}
