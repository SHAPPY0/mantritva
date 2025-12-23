package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	NONE
)

var (
	mu      sync.RWMutex
	out     io.Writer = os.Stdout
	curLv   Level     = INFO
	timeFmt           = time.RFC3339
)

// InitLogger sets global log level and optional output writer.
// levelStr: "debug","info","warn","error". If writer provided, it will be used as output.
func InitLogger(levelStr string, writer ...io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	curLv = parseLevel(levelStr)
	if len(writer) > 0 && writer[0] != nil {
		out = writer[0]
	}
}

// SetOutput replaces the logger output writer.
func SetOutput(writer io.Writer) {
	if writer == nil {
		return
	}
	mu.Lock()
	out = writer
	mu.Unlock()
}

// SetLevel sets the global logging level.
func SetLevel(levelStr string) {
	mu.Lock()
	curLv = parseLevel(levelStr)
	mu.Unlock()
}

func parseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return DEBUG
	case "warn", "warning":
		return WARN
	case "error", "err":
		return ERROR
	case "none":
		return NONE
	default:
		return INFO
	}
}

func shouldLog(l Level) bool {
	mu.RLock()
	defer mu.RUnlock()
	return l >= curLv && curLv != NONE
}

func levelString(l Level) string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Log prints structured log. fields may be nil.
func Log(l Level, msg string, _ map[string]any) {
	if !shouldLog(l) {
		return
	}
	ts := time.Now().Format("2006/01/02 - 15:04:05")
	fmt.Fprintf(out, "%s | %-5s | %s\n", ts, levelString(l), msg)
}

// Logf is like Log but with formatting.
func Logf(l Level, format string, args ...any) {
	Log(l, fmt.Sprintf(format, args...), nil)
}

// Convenience helpers:
func Debug(msg string) { Log(DEBUG, msg, nil) }
func Info(msg string)  { Log(INFO, msg, nil) }
func Warn(msg string)  { Log(WARN, msg, nil) }
func Error(msg string) { Log(ERROR, msg, nil) }

func Debugf(format string, args ...any) { Logf(DEBUG, format, args...) }
func Infof(format string, args ...any)  { Logf(INFO, format, args...) }
func Warnf(format string, args ...any)  { Logf(WARN, format, args...) }
func Errorf(format string, args ...any) { Logf(ERROR, format, args...) }

// Ensure default logger can be configured from environment if desired.
func init() {
	if lv := os.Getenv("MANTRITVA_LOG_LEVEL"); lv != "" {
		curLv = parseLevel(lv)
	}
}
