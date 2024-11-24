package debug

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

var enabled bool

// Enable sets the debug logging flag to true.
func Enable() {
	enabled = true
}

// Log logs a debug message if debug logging is enabled.
func Log(format string, v ...interface{}) {
	if enabled {
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			fn := runtime.FuncForPC(pc)
			funcName := fn.Name()
			shortFuncName := funcName[strings.LastIndex(funcName, ".")+1:]
			shortFile := filepath.Base(file)
			log.Printf("%s:%d %s: "+format, append([]interface{}{shortFile, line, shortFuncName}, v...)...)
		} else {
			log.Printf(format, v...)
		}
	}
}
