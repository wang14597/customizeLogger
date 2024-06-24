package main

import "github.com/clog/log"

func main() {
	customizeLogger := log.CustomizeLogger{}
	customizeLogger.Init()
	customizeLogger.Info("info....")
	customizeLogger.Error("error....")

	customizeLogger2 := log.CustomizeLogger{}
	customizeLogger2.SetOutPutNew(true)
	customizeLogger2.Init()
	customizeLogger2.Info("info....")
	customizeLogger2.Error("error....")

	customizeLogger2.WithFieldUID("log-1").Info("info...")
	customizeLogger2.WithFieldUID("log-1").Error("info...")
}
