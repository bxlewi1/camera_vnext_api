package config

import (
	"encoding/json"
	"log"
	"os"
)

var Instance = LoadConfig()

type Config struct {
	Port      string
	AssetPath string
}

func LoadConfig() (output Config) {
	bytes, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(bytes, &output)
	if err != nil {
		log.Fatalln(err)
	}

	return
}
