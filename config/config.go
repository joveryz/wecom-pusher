package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/TongboZhang/wecom-pusher/logger"
)

type WeComJsonConfig struct {
	Alias      string
	Receiver   string
	CorpId     string
	CorpSecret string
	AgentId    string
}

var Config []WeComJsonConfig

func LoadConfig(jsonConfigPath string) (err error) {
	file, err := ioutil.ReadFile(jsonConfigPath)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(file), &Config)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	logger.Infof("Config: %+v\n", Config)
	return err
}
