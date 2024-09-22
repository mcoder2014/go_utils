package qyapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendTextMessage(t *testing.T) {
	var webhook = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=[webhook]"
	var testMessage = "测试消息"

	err := SendTextMessage(context.Background(), webhook, testMessage)
	require.NoError(t, err)
}
