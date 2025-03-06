package debug

import (
	"fmt"
	"os"
	"time"
)

const (
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
	LogLevelError = "ERROR"
)

var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}
}

func Log(filename string, message string) {
	logLevel := GetEnvVar("LOG_LEVEL")
	if logLevel == "" {
		logLevel = LogLevelInfo
	}

	colourMap := map[string]string{
		"main":   "\033[33m", // yellow
		"pubsub": "\033[32m", // green
		"server": "\033[34m", // blue
		"p2p":    "\033[35m", // magenta
		"keys":   "\033[36m", // cyan
		"raft":   "\033[92m", // bright green
		"db":     "\033[94m", // bright blue
		"err":    "\033[91m", // bright red
		"reset":  "\033[0m",  // reset
	}

	logMessage := fmt.Sprintf("[%s] [%s] %s", filename, time.Now().Format("15:04:05"), message)

	if logLevel == LogLevelDebug || logLevel == LogLevelInfo || filename == "error" {
		if filename == "error" {
			fmt.Println(colourMap[filename] + "[" + filename + "] [" + time.Now().Format("15:04:05") + "] " + colourMap["reset"] + message)
		} else {
			fmt.Println(colourMap[filename] + "[" + filename + ".go] [" + time.Now().Format("15:04:05") + "] " + colourMap["reset"] + message)
		}
	}

	if logFile != nil {
		logFile.WriteString(logMessage + "\n")
	}
}
