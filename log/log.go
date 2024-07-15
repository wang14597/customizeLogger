package log

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	m "github.com/wang14597/customizeLogger/dataStructure/map"
	p "github.com/wang14597/customizeLogger/dataStructure/pool"
	"os"
	"strings"
	"time"
)

const DefaultTimeFormat = time.RFC3339Nano

func (c *CustomizeLogger) Init() {
	SetLogLevel(os.Getenv("LOG_LEVEL"))
	//l := logrus.StandardLogger()
	l := logrus.New() // 新对象，不实用指针对象，避免冲突
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: DefaultTimeFormat,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	})
	//l.SetReportCaller(true)
	if c.outputNew {
		channel := make(chan string, 100)
		c.LogChannel = &channel
		hook := &ServiceHookNew{
			ServiceName: "mk-agent",
			LogChannel:  c.LogChannel,
			WriterPool:  &p.WriterPool{Pool: m.New[*p.Writer]()},
		}
		hook.SetLogDir()
		c.Wg.Add(1)
		go c.startAsyncLogWriterNew(hook)
		c.ServiceHookNew = hook
		l.AddHook(hook)

		job := cron.New()
		err := job.AddFunc("*/1 * * * *", func() {
			c.ServiceHookNew.WriterPool.CleanExpiredWriter()
		})
		if err != nil {
			fmt.Println("Error adding cron job:", err)
			return
		}
		job.Start()

	}

	c.defaultLogger = l
}

func (c *CustomizeLogger) startAsyncLogWriterNew(hook *ServiceHookNew) {
	defer c.Wg.Done()
	for {
		select {
		case logEntry, ok := <-*hook.LogChannel:
			if !ok {
				return
			}
			split := strings.Split(logEntry, "@^@")
			err := writeToFileNew(split[0], split[1], hook)
			if err != nil {
				logrus.Errorf("Failed to write log to file: %v", err)
			}
		}
	}
}

func SetLogLevel(logLevel string) {
	logrus.SetLevel(getLogrusLogLevel(logLevel))
}

func getLogrusLogLevel(logLevel string) logrus.Level {
	switch logLevel {
	case "": // not set
		return logrus.InfoLevel
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	}

	panic(fmt.Sprintf("LOG_LEVEL %s is not known", logLevel))
}

func (c *CustomizeLogger) Cleanup() {
	if c.LogChannel != nil {
		close(*c.LogChannel)
	}
	if c.outputNew {
		c.ServiceHookNew.CleanPool()
	}
}

func (c *CustomizeLogger) Info(args ...interface{}) {
	c.defaultLogger.Info(args...)
}

func (c *CustomizeLogger) Error(args ...interface{}) {
	c.defaultLogger.Error(args...)
}

func (c *CustomizeLogger) WithField(key string, value interface{}) *logrus.Entry {
	return c.defaultLogger.WithField(key, value)
}

func (c *CustomizeLogger) WithFieldUID(value interface{}) *logrus.Entry {
	return c.defaultLogger.WithField("UID", value)
}
