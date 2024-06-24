package log

import (
	"github.com/sirupsen/logrus"
	p "github.com/wang14597/customizeLogger/dataStructure/pool"
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
