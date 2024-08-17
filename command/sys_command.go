package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync/atomic"

	"github.com/mcoder2014/go_utils/common"
)

type Executor struct {
	// 二进制地址
	BinaryPath string
	// 参数
	Params []string

	// 环境变量 key=val
	EnvVars []string
	Dir     string

	// 输入输出
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	// 依赖的执行器
	cmd *exec.Cmd

	// 如果执行失败，可以从这里获取错误
	Error    error
	ExitCode int
	ExitMsg  string

	done      chan struct{}
	isRunning atomic.Bool
}

func NewExecutor(binaryPath string, params ...string) *Executor {
	return &Executor{
		BinaryPath: binaryPath,
		Params:     params,
		done:       make(chan struct{}),
	}
}

// Build 参数均设置完成后调用
func (e *Executor) Build() *Executor {
	if e.cmd != nil {
		return e
	}
	e.cmd = exec.Command(e.BinaryPath, e.Params...)

	if e.Stdin != nil {
		e.cmd.Stdin = e.Stdin
	}
	if e.Stdout != nil {
		e.cmd.Stdout = e.Stdout
		e.cmd.Stderr = e.Stdout // 如未特殊指定，输出到一个地方
	}
	if e.Stderr != nil {
		e.cmd.Stderr = e.Stderr
	}
	if e.EnvVars != nil {
		e.cmd.Env = e.EnvVars
	}
	if e.Dir != "" {
		e.cmd.Dir = e.Dir
	}

	return e
}

// Exec 执行命令，阻塞式
func (e *Executor) Exec(ctx context.Context) error {
	if e.cmd == nil {
		return fmt.Errorf("executor has not been initialized")
	}

	// 异步执行命令
	go func() {
		defer close(e.done)
		defer func() {
			e.isRunning.Store(false)
		}()
		defer common.Recovery(ctx)

		e.isRunning.Store(true)
		err := e.cmd.Run() // 阻塞式
		if err != nil {
			e.Error = err
			e.ExitCode = -1 // 未知
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				e.ExitCode = exitErr.ExitCode()
				e.ExitMsg = exitErr.String()
			}
		}
	}()

	// 阻塞，等待执行结果
	select {
	case <-e.done:
		return e.Error
	}
}

func (e *Executor) Kill() error {
	if e.cmd == nil {
		return nil
	}

	return e.cmd.Process.Kill()
}

func (e *Executor) IsRunning() bool {
	if e.cmd == nil {
		return false
	}
	return e.isRunning.Load()
}
