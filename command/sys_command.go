package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

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
	errMutex sync.RWMutex
	exitErr  error
	exitCode int
	exitMsg  string

	done      chan struct{}
	sigs      chan os.Signal
	isRunning atomic.Bool
}

func (e *Executor) Error() error {
	e.errMutex.RLock()
	defer e.errMutex.RUnlock()
	return e.exitErr
}

func (e *Executor) ExitCode() int {
	e.errMutex.RLock()
	defer e.errMutex.RUnlock()
	return e.exitCode
}

func (e *Executor) ExitMsg() string {
	e.errMutex.RLock()
	defer e.errMutex.RUnlock()
	return e.exitMsg
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
	if e.cmd != nil { // 已经有数据的 cmd 可再次执行 build，会尝试终止此前的进程
		if e.IsRunning() {
			for _idx := 0; _idx < 100 && e.IsRunning(); _idx++ {
				_ = e.Kill()
				time.Sleep(1 * time.Second)
			}
		}
	}

	// 构建新的 cmd, 并指定数据
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

	// 重置程序退出的错误信息
	e.errMutex.Lock()
	defer e.errMutex.Unlock()
	e.exitMsg = ""
	e.exitCode = 0
	e.exitErr = nil

	return e
}

// Exec 执行命令，阻塞式
func (e *Executor) Exec(ctx context.Context) error {
	if e.cmd == nil {
		return fmt.Errorf("executor has not been initialized")
	}

	// 异步执行命令
	go func() {
		defer common.Recovery(ctx)
		defer func() {
			e.isRunning.Store(false)
			e.done <- struct{}{}
		}()

		e.isRunning.Store(true)
		log.Ctx(ctx).Infof("start sub program: %s param: %v", e.BinaryPath, e.Params)

		err := e.cmd.Run() // 阻塞式
		if err != nil {
			e.errMutex.Lock()
			defer e.errMutex.Unlock()
			e.exitErr = err
			e.exitCode = 255 // 未知
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				e.exitCode = exitErr.ExitCode()
				e.exitMsg = exitErr.String()
			}
			log.Ctx(ctx).WithError(err).Warnf("sub program exit with error, code: %d msg: %v", e.exitCode, e.exitMsg)
		}
	}()

	// 阻塞，等待执行结果
	select {
	case <-e.done:
		return e.exitErr
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
