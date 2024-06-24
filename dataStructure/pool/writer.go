package pool

import (
	m "github.com/clog/dataStructure/map"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

type Writer struct {
	Writer      *os.File
	ExpiresTime int // 过期时间（单位秒）
	UpdateTime  time.Time
}

func (w Writer) IsExpired() bool {
	now := time.Now()
	sub := now.Sub(w.UpdateTime)
	seconds := sub.Seconds()
	if seconds >= float64(w.ExpiresTime) {
		return true
	}
	return false
}

type WriterPool struct {
	Pool       m.ConcurrentMap[*Writer]
	WriterKeys []string
	Lock       sync.Mutex
}

func (w *WriterPool) SetIfNotExist(key string, baseDir string) error {
	_, err := w.GetIfNotExistSet(key, baseDir)
	return err
}

func (w *WriterPool) GetIfNotExistSet(key string, baseDir string) (*Writer, error) {
	writer, ok := w.Pool.Get(key)
	if !ok { // 双重锁
		w.Lock.Lock()
		_, ok := w.Pool.Get(key)
		if !ok {
			fileWriter, err := CreateFileWriter(baseDir, key)
			writer = &Writer{Writer: fileWriter, ExpiresTime: 60, UpdateTime: time.Now()}
			if err != nil {
				return nil, err
			}
			w.Pool.Set(key, writer)
			w.WriterKeys = append(w.WriterKeys, key)
		}
		w.Lock.Unlock()
	}
	return writer, nil
}

func (w *WriterPool) CleanExpiredWriter() {
	w.Lock.Lock()
	noExpiredKeys := make([]string, 0)
	for _, key := range w.WriterKeys {
		writer, _ := w.Pool.Get(key)
		if writer.IsExpired() {
			_ = writer.Writer.Close()
			w.Pool.Remove(key)
			continue
		}
		noExpiredKeys = append(noExpiredKeys, key)
	}
	w.WriterKeys = noExpiredKeys
	w.Lock.Unlock()
}

func (w *WriterPool) CleanAll() {
	w.Lock.Lock()
	for _, key := range w.WriterKeys {
		writer, _ := w.Pool.Get(key)
		_ = writer.Writer.Close()
		w.Pool.Remove(key)
	}
	w.Lock.Unlock()
}

func CreateFileWriter(logDir string, uid string) (*os.File, error) {
	currentDate := time.Now().Format("2006-01-02")
	folderPath := logDir + "/" + currentDate
	logFile := folderPath + "/" + uid + ".log"
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			logrus.Errorf("Error creating folder: %v", err)
		}
	}
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}
