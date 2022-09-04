package logger

import (
    "fmt"
	"os"
	"strings"
	"sync"
)

type LogLevel int

const (
    LOG_LEVEL_DEBUG LogLevel = iota
    LOG_LEVEL_INFO
    LOG_LEVEL_WARN
    LOG_LEVEL_ERROR
)

type Logger struct {
    level LogLevel
}

func LevelForString(s string) LogLevel {
    switch (strings.ToLower(s)) {
        case "debug":
            return LOG_LEVEL_DEBUG
        case "info":
            return LOG_LEVEL_INFO
        case "warn":
            return LOG_LEVEL_WARN
        case "error":
            return LOG_LEVEL_ERROR
    }
    return LOG_LEVEL_INFO
}

var once sync.Once
var logger *Logger

func Initialize() {
    levelString, ok := os.LookupEnv("LOG_LEVEL")

    if !ok {
        levelString = "info"
    }

    level := LevelForString(levelString)

    if logger == nil {
        once.Do(
            func() {
                logger = &Logger{level: level}
            })
    }
}

func Debug(s string) {
    if logger.level <= LOG_LEVEL_DEBUG {
        fmt.Println(s)
    }
}

func Info(s string) {
    if logger.level <= LOG_LEVEL_INFO {
        fmt.Println(s)
    }
}

func Warn(s string) {
    if logger.level <= LOG_LEVEL_WARN {
        fmt.Println(s)
    }
}

func Error(s string) {
    if logger.level <= LOG_LEVEL_ERROR {
        fmt.Println(s)
    }
}
