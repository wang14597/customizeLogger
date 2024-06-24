package log

import (
	p "github.com/clog/dataStructure/pool"
	"github.com/sirupsen/logrus"
	"sync"
)

type CustomizeLogger struct {
	defaultLogger *logrus.Logger
	Wg            sync.WaitGroup
	LogChannel    *chan string
	outputNew     bool
	*ServiceHookNew
}

func (c *CustomizeLogger) SetOutPutNew(b bool) {
	c.outputNew = b
}

type ServiceHookNew struct {
	ServiceName string
	LogChannel  *chan string
	WriterPool  *p.WriterPool
	LogDir      string
}
