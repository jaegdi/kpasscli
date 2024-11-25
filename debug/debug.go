package debug

import (
	"fmt"
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
// It includes the name of the calling function and the line number in the log message.
// Parameters:
//   - format: The format string for the log message, similar to fmt.Printf.
//   - v: The values to be formatted according to the format string.
func Log(format string, v ...interface{}) {
	for i := range v {
		val := v[i]
		valStr := fmt.Sprintf("%v", val)
		if strings.Contains(strings.ToLower(valStr), "passwor") {
			v[i] = "!Password replaced! ****"
		}
	}
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
