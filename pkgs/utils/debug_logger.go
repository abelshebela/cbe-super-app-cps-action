package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorReset  = "\033[0m"
)

func DebuggingLogger(messages ...interface{}) {
	env := os.Getenv("Go_ENV")
	isDebugEnabled := strings.ToLower(os.Getenv("ENABLE_DEBUG_LOGS")) == "true"
	timestamp := time.Now().Format(time.RFC3339)
	if (env == "dev" || env == "uat") && isDebugEnabled {
		fmt.Printf("%s[%s] ==> %s", colorGreen, timestamp, colorReset)
		fmt.Println(messages...)
	}
}

func DebuggingLoggerV2(messages []interface{}, logType ...string) {
	env := os.Getenv("NODE_ENV")
	isDebugEnabled := strings.ToLower(os.Getenv("ENABLE_DEBUG_LOGS")) == "true"
	timestamp := time.Now().Format(time.RFC3339)
	if (env == "dev" || env == "uat") && isDebugEnabled && len(messages) > 0 {
		var prefix string
		if len(logType) > 0 {
			switch strings.ToLower(logType[0]) {
			case "info":
				prefix = fmt.Sprintf("%s[%s] ==> %s", colorYellow, timestamp, colorReset)
				fmt.Print(prefix)
				fmt.Println(messages...)
				return
			case "error":
				prefix = fmt.Sprintf("%s[%s] ==> %s", colorRed, timestamp, colorReset)
				fmt.Fprint(os.Stderr, prefix)
				fmt.Fprintln(os.Stderr, messages...)
				return
			}
		}
		prefix = fmt.Sprintf("%s[%s] ==> %s", colorGreen, timestamp, colorReset)
		fmt.Print(prefix)
		fmt.Println(messages...)
	}
}
