package util

import "fmt"

// CoriPrintf displays message on standard output or log file
func CoriPrintf(format string, v ...interface{}) {
	fmt.Print(fmt.Sprintf(format, v...))
}

// CoriPrintln displays message on standard output or log file
func CoriPrintln(v ...interface{}) {
	fmt.Print(fmt.Sprintln(v...))
}
