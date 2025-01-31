package sendcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SendCloudAPI interface {
	SendPlainEmail(ctx context.Context, param *NormalPlanEmailParam) (*BaseResp, error)
}

var defaultClient SendCloudAPI

type SendCloudAPIImpl struct {
	apiUser string
	apiKey  string
}

func InitSendCloudAPI(apiUser string, apiKey string) {
	impl := &SendCloudAPIImpl{
		apiUser: apiUser,
		apiKey:  apiKey,
	}
	defaultClient = impl
}

func NewSendCloudAPI(apiUser string, apiKey string) SendCloudAPI {
	return &SendCloudAPIImpl{
		apiUser: apiUser,
		apiKey:  apiKey,
	}
}

func SetDefaultSendCloudAPIAPI(apiImpl SendCloudAPI) {
	defaultClient = apiImpl
}

func GetDefault() SendCloudAPI {
	return defaultClient
}

func (c *SendCloudAPIImpl) SendPlainEmail(ctx context.Context, param *NormalPlanEmailParam) (*BaseResp, error) {
	nParam := param.ToNormalSendEmailParam(c.apiUser, c.apiKey)
	return c.NormalSendEmail(ctx, nParam)
}

func (c *SendCloudAPIImpl) NormalSendEmail(ctx context.Context, param *NormalSendEmailParam) (*BaseResp, error) {
	const url = "https://api.sendcloud.net/apiv2/mail/send"

	body, err := json.Marshal(param)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	urlValues := param.ToURLValue()
	postBody := bytes.NewBufferString(urlValues.Encode())

	httpReq, err := http.NewRequest(http.MethodPost, url, postBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http client request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http response body failed, err :%w", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("send feishu message failed, status code: %d, body: %v", httpResp.StatusCode, string(body))
	}
	var resp BaseResp
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &resp, nil
}
