package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	// log.Lshortfile 支持显示文件名和代码行号
	// 颜色为红色
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	// 颜色为蓝色
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

// log methods
var (
	Error = errorLog.Println
	Errorf = errorLog.Printf
	Info = errorLog.Println
	Infof = errorLog.Printf
)

// log levels
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// SetLevel controls log level
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}
}
