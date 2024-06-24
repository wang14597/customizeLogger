package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func (h *ServiceHookNew) SetLogDir() {
	logDir := os.Getenv("LOG_PATH")
	if logDir == "" {
		logDir, _ = os.Getwd()
	}
	logDir = logDir + "/logs"
	h.LogDir = logDir
}

func (h *ServiceHookNew) CleanPool() {
	h.WriterPool.CleanAll()
}

func (h *ServiceHookNew) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *ServiceHookNew) Fire(entry *logrus.Entry) error {
	entry.Data["service"] = h.ServiceName
	message := createMessage(entry)
	uid := strings.Split(message, "@^@")[0]
	err := h.WriterPool.SetIfNotExist(uid, h.LogDir)
	if err != nil {
		return err
	}
	*h.LogChannel <- fmt.Sprintf(message)
	return nil
}
