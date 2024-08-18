package custom_bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BaseElement struct {
	Tag string `json:"tag"`
}

type FeishuCardMessage struct {
	MsgType string       `json:"msg_type"`
	Card    *CardContent `json:"card"`
}

type CardContent struct {
	Header   *CardHeader            `json:"header"`
	Elements []interface{}          `json:"elements"`
	Config   map[string]interface{} `json:"config,omitempty"`
}

type CardHeader struct {
	Title    *Text  `json:"title"`
	Template string `json:"template,omitempty"`
}

type MarkdownElement struct {
	BaseElement
	Text *Text `json:"text"`
}

type Text struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

const (
	MsgTypeInteractive = "interactive"

	TemplateRed    = "red"
	TemplateGreen  = "green"
	TemplateOrange = "orange"

	ElementTagDiv       = "div"
	ElementTagColumnSet = "column_set"
	ElementTagMarkdown  = "markdown"

	TextTagMD        = "lark_md"
	TextTagPlainText = "plain_text"

	FlexModeFlow = "flow"

	BackgroundDefault = "default"

	VerticalAlignTop = "top"

	TextAlignCenter = "center"
)

func SendFeishuMessage(ctx context.Context, url string, message string) error {
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(message)))
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

func SendErrorMessage(ctx context.Context, url string, title string, msg string, args ...any) error {
	cardMessage := &FeishuCardMessage{
		MsgType: MsgTypeInteractive,
		Card: &CardContent{
			Header: &CardHeader{
				Template: TemplateRed,
				Title: &Text{
					Tag:     TextTagPlainText,
					Content: title, // 标题
				},
			},
			Elements: []interface{}{
				&MarkdownElement{
					BaseElement: BaseElement{
						Tag: ElementTagDiv,
					},
					Text: &Text{
						Tag:     TextTagPlainText,
						Content: fmt.Sprintf(msg, args...), // 内容
					},
				},
			},
			Config: map[string]interface{}{
				"wide_screen_mode": true,
			},
		},
	}
	js, _ := json.MarshalIndent(cardMessage, "", " ")
	return SendFeishuMessage(ctx, url, string(js))
}
