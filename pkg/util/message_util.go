package util

import (
	"fmt"
	"time"
)

// Get current time string
func GetCurTimeStr() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05.000")
}

// CoriPrintf displays message on standard output or log file
func CoriPrintf(format string, v ...interface{}) {
	format = fmt.Sprintf("[ %s ] %s", GetCurTimeStr(), format)
	fmt.Print(fmt.Sprintf(format, v...))
}

// CoriPrintln displays message on standard output or log file
func CoriPrintln(v ...interface{}) {
	fmt.Print(fmt.Sprintf("[ %s ] ", GetCurTimeStr()))
	fmt.Print(fmt.Sprintln(v...))
}
