package chai

import (
	"fmt"
	"syscall/js"
)

var logs_num int

func LogF(_format_string string, args ...interface{}) {
	// fmt.Printf(_format_string+"\n", args...)
	js.Global().Get("console").Call("log", fmt.Sprintf(_format_string, args...))

	//consoleLogFormattedString(_format_string, "", args...)
	logs_num++
}

func WarningF(_format_string string, args ...interface{}) {
	// fmt.Printf(_format_string+"\n", args...)
	js.Global().Get("console").Call("warn", fmt.Sprintf(_format_string, args...))
	//consoleLogFormattedString(_format_string, "style='color:#DBD51D'", args...)
	logs_num++
}

func ErrorF(_format_string string, args ...interface{}) {
	// fmt.Printf(_format_string+"\n", args...)
	js.Global().Get("console").Call("error", fmt.Sprintf(_format_string, args...))
	//consoleLogFormattedString(_format_string, "style='color:red'", args...)
	logs_num++
}

// Print out a message if condition is false
func Assert(_condition bool, _error_message string, args ...interface{}) {
	if !_condition {
		ErrorF(_error_message, args...)
		panic("PROGRAM PANICKED")
	}
}

func AssertNot(_condition bool, _error_message string, args ...interface{}) {
	if _condition {
		ErrorF(_error_message, args...)
		panic("PROGRAM PANICKED")
	}
}

// func consoleLogFormattedString(_format_string string, _extra_html_args string, args ...interface{}) {
// 	debug_console.Set("innerHTML", debug_console.Get("innerHTML").String()+"<p class='debug-line' "+_extra_html_args+"><strong>["+time.Now().Format("2006-01-02 15:04:05")+"]:</strong> "+fmt.Sprintf(_format_string, args...)+"</p>")
// }

// func consoleLogString(_format_string string, _extra_html_args string) {
// 	debug_console.Set("innerHTML", debug_console.Get("innerHTML").String()+"<p class='debug-line' "+_extra_html_args+"><strong>["+time.Now().Format("2006-01-02 15:04:05")+"]:</strong> "+_format_string+"</p>")
// }

func terminateMessages() {
	js.Global().Get("postMessage").Invoke("terminate")
}
