package chai

import (
	"fmt"
)

func LogF(_format_string string, args ...interface{}) {
	fmt.Printf(_format_string+"\n", args...)
}

// Print out a message if condition is false
func Assert(_condition bool, _error_message string, args ...interface{}) {
	if !_condition {
		LogF(_error_message, args...)
		panic("PROGRAM PANICKED")
	}
}

func AssertNot(_condition bool, _error_message string, args ...interface{}) {
	if _condition {
		LogF(_error_message, args...)
		panic("PROGRAM PANICKED")
	}
}
