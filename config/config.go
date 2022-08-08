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

var ConfigMap map[string]WeComJsonConfig

func LoadConfig(jsonConfigPath string) (err error) {
	file, err := ioutil.ReadFile(jsonConfigPath)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	var configs []WeComJsonConfig
	err = json.Unmarshal([]byte(file), &configs)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	ConfigMap = make(map[string]WeComJsonConfig)
	for _, s := range configs {
		ConfigMap[s.Alias] = s
	}

	logger.Infof("ConfigMap: %+v\n", ConfigMap)
	return err
}
