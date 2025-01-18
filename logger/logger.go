package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"sgithub.com/GGBusuioc/go-cache/config"
)

type Logger struct {
	logger *log.Logger
	level  config.LogLevel
}

func NewLogger(config *config.Config) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", 0),
		level:  config.LogLevel,
	}
}

func (l *Logger) formatMessage(level config.LogLevel, message string) string {
	levelMessage := ""
	switch level {
	case config.DEBUG:
		levelMessage = "DEBUG"
	case config.INFO:
		levelMessage = "INFO"
	case config.ERROR:
		levelMessage = "ERROR"
	}

	timestamp := time.Now().Format("2025-01-06 15:00:00")
	return fmt.Sprintf("[%s] [%s] %s", timestamp, levelMessage, message)
}

func (l *Logger) shouldLog(levelMessage config.LogLevel) bool {
	return levelMessage >= l.level
}

func (l *Logger) Debug(message string) {
	if l.shouldLog(config.DEBUG) {
		l.logger.Println(l.formatMessage(config.DEBUG, message))
	}
}

func (l *Logger) Info(message string) {
	if l.shouldLog(config.INFO) {
		l.logger.Println(l.formatMessage(config.INFO, message))
	}
}

func (l *Logger) Error(message string) {
	if l.shouldLog(config.ERROR) {
		l.logger.Println(l.formatMessage(config.ERROR, message))
	}
}
