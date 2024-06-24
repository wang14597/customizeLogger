package log

import (
	"fmt"
	"testing"
	"time"
)

func add(a int, b int) int {
	return a + b
}

type MyAddType interface {
	// int | string
	any
}

func getTime[T MyAddType](f func(T, T) T) func(T, T) T {
	return func(params1 T, params2 T) T {
		t1 := time.Now()
		res := f(params1, params1)
		t2 := time.Now()
		fmt.Println("计算时间为：", t2.Sub(t1))
		return res
	}
}

func TestGetTime(t *testing.T) {
	f := getTime(add)
	println(f(2, 1))
}

func TestLogInit(t *testing.T) {
	customizeLogger := CustomizeLogger{}
	customizeLogger.Init()
	for i := 0; i < 100000; i++ {
		customizeLogger.defaultLogger.Error(fmt.Sprintf("error: %d", i))
	}
}

func TestLogWithIdNew(t *testing.T) {
	start := time.Now()
	customizeLogger := CustomizeLogger{}
	customizeLogger.SetOutPutNew(true)
	customizeLogger.Init()
	for i := 0; i < 10000; i++ {
		customizeLogger.defaultLogger.WithField("UID", "api-hof-51").Error(fmt.Sprintf("error: %d", i))
	}
	for i := 0; i < 10000; i++ {
		customizeLogger.defaultLogger.WithField("UID", "api-hof-41").Error(fmt.Sprintf("error: %d", i))
	}
	for i := 0; i < 10000; i++ {
		customizeLogger.defaultLogger.WithField("UID", "api-hof-61").Error(fmt.Sprintf("error: %d", i))
	}
	for i := 0; i < 10000; i++ {
		customizeLogger.defaultLogger.WithField("UID", "api-hof-71").Error(fmt.Sprintf("error: %d", i))
	}
	end := time.Now()
	customizeLogger.Cleanup()
	customizeLogger.Wg.Wait()
	duration := end.Sub(start)
	fmt.Println("时间差:", duration)
}

func TestLogWithIdNewClean(t *testing.T) {
	customizeLogger := CustomizeLogger{}
	customizeLogger.SetOutPutNew(true)
	customizeLogger.Init()
	for i := 0; i < 10; i++ {
		customizeLogger.defaultLogger.WithField("UID", "api-hof-1").Error(fmt.Sprintf("error: %d", i))
	}
	for i := 0; i < 10; i++ {
		customizeLogger.defaultLogger.WithField("UID", "api-hof-2").Error(fmt.Sprintf("error: %d", i))
	}
	for i := 0; i < 10; i++ {
		customizeLogger.defaultLogger.WithField("UID", "api-hof-3").Error(fmt.Sprintf("error: %d", i))
	}
	for i := 0; i < 10; i++ {
		customizeLogger.defaultLogger.WithField("UID", "api-hof-4").Error(fmt.Sprintf("error: %d", i))
	}
	customizeLogger.Cleanup()
	customizeLogger.Wg.Wait()
}
