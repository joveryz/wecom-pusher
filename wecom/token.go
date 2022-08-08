package wecom

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/logger"
)

type tokenType struct {
	Token      string
	Expiration time.Time
}

var tokenMap map[string]tokenType

func GetToken(alias string) (tokenString string, err error) {
	_, ok := config.ConfigMap[alias]
	if !ok {
		return "", errors.New("alias not found in configs: " + alias)
	}

	token, ok := tokenMap[alias]
	if ok && !token.Expiration.IsZero() && time.Now().Before(token.Expiration) {
		return token.Token, err
	}

	logger.Infof("No valid token for %s, will require a new one.", alias)

	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?" + "corpid=" + config.ConfigMap[alias].CorpId + "&corpsecret=" + config.ConfigMap[alias].CorpSecret
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(data, &objmap)
	if err != nil {
		return "", err
	}

	tokenString = string(*objmap["access_token"])
	tokenString = strings.Replace(tokenString, "\"", "", -1)

	newToken := tokenType{
		Token:      tokenString,
		Expiration: time.Now().Add(1 * time.Hour),
	}
	tokenMap[alias] = newToken
	return tokenMap[alias].Token, err
}

func init() {
	tokenMap = make(map[string]tokenType)
}
