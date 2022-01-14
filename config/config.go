package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type Config struct {
	ConnectionString string `json:"connectionString"`
	EmailPassword    string `json:"emailPassword"`
}

var (
	singleton = Config{}
	once      = sync.Once{}
)

func GetConfig(path string) Config {
	once.Do(func() {
		jsonFile, _ := os.Open(path)
		defer jsonFile.Close()

		jsonData, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(jsonData, &singleton)
	})

	return singleton
}
