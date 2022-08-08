package config

import (
	"encoding/base64"
	"encoding/json"
	"github.com/TongboZhang/wecom-pusher/logger"
	"io/ioutil"
	"os"
)

type WeComJsonConfig struct {
	Alias      string
	Receiver   string
	CorpId     string
	CorpSecret string
	AgentId    string
}

type ConfigArray struct {
	User         string
	Password     string
	WeComConfigs []WeComJsonConfig
}

type ConfigMap struct {
	User         string
	Password     string
	Token        string
	Aliases      []string
	WeComConfigs map[string]WeComJsonConfig
}

var Config ConfigMap

func LoadConfig(jsonConfigPath string) (err error) {
	file, err := ioutil.ReadFile(jsonConfigPath)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	var configArray ConfigArray
	err = json.Unmarshal([]byte(file), &configArray)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	Config.User = configArray.User
	Config.Password = configArray.Password
	Config.Token = base64.StdEncoding.EncodeToString([]byte(configArray.User + ":" + configArray.Password))
	Config.WeComConfigs = make(map[string]WeComJsonConfig)
	for _, s := range configArray.WeComConfigs {
		Config.WeComConfigs[s.Alias] = s
		Config.Aliases = append(Config.Aliases, s.Alias)
	}

	logger.Infof("Config: %+v\n", Config)
	return err
}
