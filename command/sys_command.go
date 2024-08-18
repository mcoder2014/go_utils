package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/mcoder2014/go_utils/common"
	"github.com/mcoder2014/go_utils/log"
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

	// 额外信息，供使用方自行使用
	Extra map[string]string

	// 依赖的执行器
	cmd *exec.Cmd

	// 如果执行失败，可以从这里获取错误
	Error    error
	ExitCode int
	ExitMsg  string

	done      chan struct{}
	sigs      chan os.Signal
	isRunning atomic.Bool
}

func NewExecutor(binaryPath string, params ...string) *Executor {
	e := &Executor{
		BinaryPath: binaryPath,
		Params:     params,
		done:       make(chan struct{}),
		sigs:       make(chan os.Signal),
		Extra:      make(map[string]string),
	}
	signal.Notify(e.sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	return e
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
		defer func() {
			e.done <- struct{}{}
		}()
		defer func() {
			e.isRunning.Store(false)
		}()
		defer common.Recovery(ctx)

		e.isRunning.Store(true)
		log.Ctx(ctx).Infof("start sub program: %s param: %v", e.BinaryPath, e.Params)
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
	case sig := <-e.sigs:
		log.Ctx(ctx).Warnf("received signal: %s, kill sub program", sig.String())
		_ = e.Kill()
		return fmt.Errorf("process killed by signal: %s", sig.String())
	}
}

func (e *Executor) Kill() error {
	if e.cmd == nil {
		return nil
	}
	log.Ctx(context.Background()).Infof("kill sub program: %s pid=%v param=%v",
		e.BinaryPath, e.cmd.Process.Pid, e.Params)
	return e.cmd.Process.Kill()
}

func (e *Executor) IsRunning() bool {
	if e.cmd == nil {
		return false
	}
	return e.isRunning.Load()
}
