package sendcloud

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendCloudAPIImpl_SendPlainEmail(t *testing.T) {

	client := NewSendCloudAPI("test", "apikey")

	resp, err := client.SendPlainEmail(context.Background(), &NormalPlanEmailParam{
		From:     "sys@mcoder.cc",
		FromName: "system",
		To:       "mcoder2014@sina.com",
		Subject:  "测试邮件 api 2",
		Plain:    "这是一个段内容，用于测试邮件 API 是否通畅",
	})
	require.NoError(t, err)
	t.Logf("resp: %+v", resp)
	content, err := json.MarshalIndent(resp, "", "\t")
	require.NoError(t, err)
	t.Logf("format:%v", string(content))

}
