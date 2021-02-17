package logger

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func getFileAndLine() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		panic("Could not get context info for logger!")
	}
	return file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
}

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString, fmt.Sprintln(time.Now().Format("02-01-2006 15:04:05")+": "+getFileAndLine()+": "+fmt.Sprint(args...)))
	}
	return sprint
}

func L(color func(...interface{}) string, message string, param ...interface{}) {
	log.Println(color(fmt.Sprintf(message, param...)))
}
