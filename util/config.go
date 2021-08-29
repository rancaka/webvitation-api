package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Config struct {
	Environment string
	Firebase    struct {
		Credentials string
	}
	Storage struct {
		DB struct {
			DataSourceName string
		}
		Bucket struct {
			Name string
		}
	}
}

func (c Config) IsProduction() bool {
	return c.Environment == "production"
}

var (
	config     Config
	configOnce sync.Once
)

func getConfig() Config {
	configOnce.Do(func() {

		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatalln("[BOOTING] error when trying to get working directory")
			return
		}

		path := workingDir + "/config/config.json"
		b, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("[BOOTING] error when trying to read config: %v\n", err)
			return
		}

		err = json.Unmarshal(b, &config)
		if err == nil {
			fmt.Printf("[BOOTING] successfully initialize config from %v!\n", path)
			return
		}
	})

	return config
}
