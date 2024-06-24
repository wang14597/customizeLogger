package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

func createMessage(entry *logrus.Entry) string {
	m := make(map[string]interface{})
	m["level"] = entry.Level
	m["timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	m["message"] = entry.Message
	//m["func"] = entry.Caller.Function + " line:" + fmt.Sprint(entry.Caller.Line)
	data := entry.Data
	for k, v := range data {
		if k == "UID" {
			continue
		}
		m[k] = v
	}
	value, exists := entry.Data["UID"]
	s, _ := entry.String()
	if exists {
		return fmt.Sprintf("%s@^@%s", value, s)
	}
	return fmt.Sprintf("%s@^@%s", "log", s)
}

func writeToFileNew(uid string, msg string, hook *ServiceHookNew) error {
	writer, err := hook.WriterPool.GetIfNotExistSet(uid, hook.LogDir)
	if err != nil {
		logrus.Errorf("Error get file writer: %v", err)
	}
	_, err = writer.Writer.WriteString(msg)
	writer.UpdateTime = time.Now()
	if err != nil {
		return err
	}
	return nil
}
