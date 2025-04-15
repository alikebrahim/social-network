package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

type stdLogger struct {
	config    Config
	fields    map[string]interface{}
	pkg       string
	mu        sync.Mutex
	fileLog   *log.Logger // Logger for file output
	console   bool        // Whether to output to console
	useColors bool        // Whether to use colors
	logFile   *os.File    // File handle for log file
}

var defaultLogger Logger

func init() {
	defaultConfig := Config{
		Level:     INFO,
		ToConsole: true,
		UseColors: true,
	}
	defaultLogger = newStdLogger(defaultConfig)
}

func getLevelColor(level Level) string {
	switch level {
	case DEBUG:
		return colorCyan
	case INFO:
		return colorGreen
	case WARN:
		return colorYellow
	case ERROR:
		return colorRed
	default:
		return colorReset
	}
}

func newStdLogger(config Config) *stdLogger {
	var fileLog *log.Logger
	var logFile *os.File
	
	// Setup file logging if needed
	if config.LogFile != "" {
		var err error
		logFile, err = os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			fileLog = log.New(logFile, "", log.LstdFlags)
		}
	}

	return &stdLogger{
		config:    config,
		fields:    make(map[string]interface{}),
		fileLog:   fileLog,
		console:   config.ToConsole,
		useColors: config.UseColors,
		logFile:   logFile,
	}
}

func (l *stdLogger) log(level Level, msg string, keyvals ...interface{}) {
	if level < l.config.Level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	fields := make(map[string]interface{})
	for k, v := range l.fields {
		fields[k] = v
	}

	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			fields[fmt.Sprint(keyvals[i])] = keyvals[i+1]
		}
	}

	_, file, line, _ := runtime.Caller(2)
	shortFile := filepath.Base(file)
	pkgFile := l.pkg
	if pkgFile == "" {
		pkgFile = filepath.Dir(file)
		pkgFile = filepath.Base(pkgFile)
	}
	location := fmt.Sprintf("%s/%s:%d", pkgFile, shortFile, line)

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	var fieldsStr string
	if len(fields) > 0 {
		pairs := make([]string, 0, len(fields))
		for k, v := range fields {
			pairs = append(pairs, fmt.Sprintf("%s=%v", k, v))
		}
		fieldsStr = " {" + strings.Join(pairs, ", ") + "}"
	}

	// Log to console if enabled
	if l.console {
		var levelStr string
		if l.useColors {
			// Use color for the level
			levelColor := getLevelColor(level)
			levelStr = fmt.Sprintf("%s%s%s%s", levelColor, colorBold, level.String(), colorReset)
		} else {
			levelStr = level.String()
		}
		
		// Format and print to console
		logLine := fmt.Sprintf("%s [%s] [%s] %s%s\n", 
			timestamp, 
			levelStr, 
			location, 
			msg, 
			fieldsStr)
		
		fmt.Print(logLine)
	}
	
	// Log to file if configured
	if l.fileLog != nil {
		// Plain text format for file
		plainLogLine := fmt.Sprintf("[%s] [%s] %s%s", 
			level.String(), 
			location, 
			msg, 
			fieldsStr)
		l.fileLog.Println(plainLogLine)
	}
}

func (l *stdLogger) Debug(msg string, keyvals ...interface{}) {
	l.log(DEBUG, msg, keyvals...)
}

func (l *stdLogger) Info(msg string, keyvals ...interface{}) {
	l.log(INFO, msg, keyvals...)
}

func (l *stdLogger) Warn(msg string, keyvals ...interface{}) {
	l.log(WARN, msg, keyvals...)
}

func (l *stdLogger) Error(msg string, keyvals ...interface{}) {
	l.log(ERROR, msg, keyvals...)
}

func (l *stdLogger) WithField(key string, value interface{}) Logger {
	newLogger := &stdLogger{
		config:    l.config,
		fileLog:   l.fileLog,
		console:   l.console,
		useColors: l.useColors,
		logFile:   l.logFile,
		pkg:       l.pkg,
		fields:    make(map[string]interface{}),
	}

	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	newLogger.fields[key] = value

	return newLogger
}

func (l *stdLogger) WithFields(fields map[string]interface{}) Logger {
	newLogger := &stdLogger{
		config:    l.config,
		fileLog:   l.fileLog,
		console:   l.console,
		useColors: l.useColors,
		logFile:   l.logFile,
		pkg:       l.pkg,
		fields:    make(map[string]interface{}),
	}

	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

func (l *stdLogger) WithPackage(pkg string) Logger {
	newLogger := &stdLogger{
		config:    l.config,
		fileLog:   l.fileLog,
		console:   l.console,
		useColors: l.useColors,
		logFile:   l.logFile,
		pkg:       pkg,
		fields:    make(map[string]interface{}),
	}

	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// Close closes any open file handles
func (l *stdLogger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}