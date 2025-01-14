package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"sync/atomic"
)

const (
	LogError LogLevel = C.LOG_ERROR
	LogWarn  LogLevel = C.LOG_WARN
	LogInfo  LogLevel = C.LOG_INFO
	LogDebug LogLevel = C.LOG_DEBUG
	LogTrace LogLevel = C.LOG_TRACE
)

type LogLevel uint32
type LogFunc func(level LogLevel, message string)

var logFunc atomic.Value

func SetLogFunc(fn LogFunc) {
	logFunc.Store(fn)
}

func logger() LogFunc {
	fn := logFunc.Load()
	if fn == nil {
		return defaultLogger
	}

	return fn.(LogFunc)
}

func defaultLogger(level LogLevel, message string) {
	switch level {
	case C.LOG_ERROR:
		fmt.Printf("[ERROR] %s\n", message)
	case C.LOG_WARN:
		fmt.Printf("[WARN] %s\n", message)
	case C.LOG_INFO:
		fmt.Printf("[INFO] %s\n", message)
	case C.LOG_DEBUG:
		fmt.Printf("[DEBUG] %s\n", message)
	case C.LOG_TRACE:
		fmt.Printf("[TRACE] %s\n", message)
	}
}
