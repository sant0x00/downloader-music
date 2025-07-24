package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sant0x00/downloader-music/internal/domain"
)

type SimpleLogger struct {
	level      string
	outputFile string
	logger     *log.Logger
}

func NewSimpleLogger(level, outputFile string) (domain.Logger, error) {
	logger := &SimpleLogger{
		level:      level,
		outputFile: outputFile,
	}

	var output *os.File
	if outputFile != "" {
		file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		output = file
	} else {
		output = os.Stdout
	}

	logger.logger = log.New(output, "", 0)
	return logger, nil
}

func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
	if l.shouldLog("info") {
		l.writeLog("INFO", msg, fields...)
	}
}

func (l *SimpleLogger) Error(msg string, err error, fields ...interface{}) {
	if l.shouldLog("error") {
		allFields := append(fields, "error", err.Error())
		l.writeLog("ERROR", msg, allFields...)
	}
}

func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
	if l.shouldLog("debug") {
		l.writeLog("DEBUG", msg, fields...)
	}
}

func (l *SimpleLogger) Warn(msg string, fields ...interface{}) {
	if l.shouldLog("warn") {
		l.writeLog("WARN", msg, fields...)
	}
}

func (l *SimpleLogger) shouldLog(level string) bool {
	levels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
	}

	currentLevel, exists := levels[l.level]
	if !exists {
		currentLevel = 1 // info como padrão
	}

	logLevel, exists := levels[level]
	if !exists {
		return false
	}

	return logLevel >= currentLevel
}

func (l *SimpleLogger) writeLog(level, msg string, fields ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s] %s: %s", timestamp, level, msg)

	if len(fields) > 0 {
		logMsg += " |"
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				logMsg += fmt.Sprintf(" %v=%v", fields[i], fields[i+1])
			}
		}
	}

	l.logger.Println(logMsg)
}
