package wecom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/logger"
)

func SendTextMessage(content string, alias string) (err error) {
	logger.Infof("send text message: %s", content)
	token, err := GetToken(alias)
	if err != nil {
		logger.Errorf("get token failed for %s, error: %v", alias, err)
		return err
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
	msg := TextMessage{
		Touser:  config.ConfigMap[alias].Receiver,
		Msgtype: "text",
		Agentid: config.ConfigMap[alias].AgentId,
	}

	msg.Text.Content = content

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("send text message failed, error: %v", err)
		return err
	}

	_, err = post(url, data)
	if err != nil {
		logger.Errorf("send text message failed, error: %v", err)
		return err
	}

	logger.Infof("send text message succeeded: %s", content)
	return err
}

func SendTextCardMessage(content string, title string, cardUrl string, alias string) (err error) {
	logger.Infof("send text card message: %s", content)
	token, err := GetToken(alias)
	if err != nil {
		logger.Errorf("get token failed for %s, error: %+v", alias, err)
		return err
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)

	msg := TextCardMessage{
		Touser:  config.ConfigMap[alias].Receiver,
		Msgtype: "textcard",
		Agentid: config.ConfigMap[alias].AgentId,
	}

	msg.TextCard.Title = title
	msg.TextCard.URL = cardUrl
	msg.TextCard.Description = content

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("send text card message failed, error: %+v", err)
		return err
	}

	_, err = post(url, data)
	if err != nil {
		logger.Errorf("send text card message failed, error: %+v", err)
		return err
	}

	logger.Infof("send text card message succeeded: %s", content)
	return err
}

func post(url string, requestBody []byte) (responseBody []byte, err error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, _ = ioutil.ReadAll(resp.Body)
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(responseBody, &objmap)
	if err != nil {
		return nil, err
	}

	errcode := string(*objmap["errcode"])
	if errcode != "0" {
		return nil, errors.New("post json to wecom failed " + string(*objmap["errmsg"]))
	}
	return responseBody, err
}
