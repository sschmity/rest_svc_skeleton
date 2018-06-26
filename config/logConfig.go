package config

import (
	"sync"
	"os"
	"log"
)

type LogConfigSingleton struct {
	filePath string
	fileName string
}

const defaultFilePath= "/log"
const defaultFileName= "server.log"

var logConfigInstance *LogConfigSingleton
var onceLog sync.Once

func GetlogConfigInstance() *LogConfigSingleton {
	onceLog.Do(func() {
		logConfigInstance = &LogConfigSingleton{filePath: defaultFilePath, fileName: defaultFileName}
	})
	return logConfigInstance
}

func SetLogger() (*os.File) {
	//create log file
	f, err := os.OpenFile(logConfigInstance.GetAbsLogPath(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()
	//set output of logs to f
	log.SetOutput(f)
	//test case
	log.Println("Log Initialization Started")
	return f
}

func (LogConfigSingleton) GetAbsLogPath() string {
	return logConfigInstance.filePath + "\\" + logConfigInstance.fileName
}

func (LogConfigSingleton) SetAbsLogPath(aFilePath string, aFileName string) {
	logConfigInstance.fileName = aFileName
	logConfigInstance.filePath = aFilePath
}
