package command

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSystemCommand(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	t.Logf("working dir: %s", wd)

	executor := NewExecutor("../output/bin/test_web_server")
	executor.Stdout = os.Stdout
	executor.Build()

	go func() {
		err := executor.Exec(context.Background())
		//require.NoError(t, err)
		t.Logf("program exited: %s", err)
	}()

	time.Sleep(2 * time.Second)
	t.Logf("waiting for command execute")
	require.True(t, executor.IsRunning())

	time.Sleep(5 * time.Second)
	t.Logf("kill program")
	err = executor.Kill()
	require.NoError(t, err)

	time.Sleep(2 * time.Second)
	t.Logf("program exited: %s exited code:%v exited msg:%v", err, executor.ExitCode, executor.exitMsg)
	require.False(t, executor.IsRunning())
}
