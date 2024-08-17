package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LogIDKey       = "log_id"
	defaultLogPath = "/var/log/home_server/run.log"
)

var defaultLogOutput io.Writer
var defaultLogOutputOnce sync.Once

func Ctx(ctx context.Context) *logrus.Entry {
	logger := logrus.StandardLogger()
	fields := logrus.Fields{}
	// logID
	if c := ctx.Value(LogIDKey); c != nil {
		fields[LogIDKey] = c
	}

	return logger.WithFields(fields)
}

type MyLogConfig struct {
	SavePath string `json:"save_path"`
}

func Init(configs ...*MyLogConfig) error {
	if defaultLogOutput != nil {
		return fmt.Errorf("the default log is already init, please do not init again")
	}
	var config = &MyLogConfig{
		SavePath: defaultLogPath,
	}
	if len(configs) > 0 {
		config = configs[0]
	}

	// logrus init
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	output := getDefaultOutput(config.SavePath)
	logrus.SetOutput(output)
	return nil
}

func getDefaultOutput(filepath string) io.Writer {
	defaultLogOutputOnce.Do(func() {
		l := &lumberjack.Logger{
			Filename:   filepath,
			MaxSize:    100, // megabytes
			MaxBackups: 64,
			MaxAge:     15,    //days
			Compress:   false, // disabled by default
		}

		myLogWriter := &MyLogWriter{
			Logger: l,
			ToStd:  true,
		}
		defaultLogOutput = myLogWriter

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP)

		go func() {
			for {
				<-c
				err := l.Rotate()
				if err != nil {
					logrus.Errorf("log rotate error: %v", err)
				}
			}
		}()
	})
	return defaultLogOutput
}
