package qyapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TextContent struct {
	Content string `json:"content"`
}

type MessageStruct struct {
	MsgType string       `json:"msgtype"`
	Text    *TextContent `json:"text,omitempty"`
}

const (
	MsgTypeText = "text"
)

func SendTextMessage(ctx context.Context, webhook string, content string) error {
	msg := &MessageStruct{
		MsgType: MsgTypeText,
		Text: &TextContent{
			Content: content,
		},
	}
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, webhook, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("read http response body failed, err :%w", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("send feishu message failed, status code: %d, body: %v", httpResp.StatusCode, string(body))
	}
	return nil

}
