package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	accessToken = os.Getenv("ACCESS_TOKEN")
	feishuURL   = "https://open.feishu.cn/open-apis/bot/v2/hook"
)

type FeishuMsg struct {
	MsgType string `json:"msg_type"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	reqBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %s", err)
		http.Error(w, "ioutil.ReadAll: "+err.Error(), http.StatusInternalServerError)
		return
	}

	msg := FeishuMsg{
		MsgType: "text",
		Text: struct {
			Content string "json:\"content\""
		}{
			Content: string(reqBytes),
		},
	}

	respData, err := sendFeishuMsg(accessToken, &msg)
	if err != nil {
		log.Printf("sendFeishuMsg: %s", err)
		http.Error(w, "sendFeishuMsg: "+err.Error(), http.StatusInternalServerError)
	} else {
		b, _ := json.Marshal(msg)
		log.Printf("sendFeishuMsg: %s", b)
		w.Write(respData)
	}
}

func sendFeishuMsg(accessToken string, msg *FeishuMsg) (respData []byte, err error) {
	rawData, _ := json.Marshal(msg)
	buf := bytes.NewBuffer(rawData)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", feishuURL, accessToken), buf)
	req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("sendFeishuMsg: %s", err)
	}
	return
}
