/*
   peter-bird.com/logger

   Usage:

	const (
			APP_NAME        = "A1-C0D3R"
			APP_SERVICE     = "MAIN"
			NEW_LOG_ERR_FMT = "Error initializing logger: %s\n"
		)
		log, err := logger.New(
			logger.Info,
			fmt.Sprintf("%s %s", APP_NAME, APP_SERVICE),
			cfg.Logger.LogFile,
		)
		if err != nil {
			log.Fatalf(NEW_LOG_ERR_FMT, err)
			return
		}
*/

package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	OpenLogErrFmt = "Failed to open log file: %s"

	DebugPrefix = " DEBUG: "
	InfoPrefix  = " INFO : "
	WarnPrefix  = " WARN : "
	ErrorPrefix = " ERROR: "

	FileModeRW = 0666
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

// Logger defines the interface for logging
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// CustomLogger implements the Logger interface
type CustomLogger struct {
	logger   *log.Logger
	logLevel LogLevel
	name     string
}

// New creates a new CustomLogger. If the file path is provided, it attempts to use it as the log output.
func New(logLevel LogLevel, name, filePath string) (*CustomLogger, error) {
	var output *os.File
	var err error

	if filePath != "" {
		output, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, FileModeRW)
		if err != nil {
			return nil, fmt.Errorf(OpenLogErrFmt, err)
		}
	} else {
		output = os.Stdout
	}

	return &CustomLogger{
		logger:   log.New(output, "", log.Ldate|log.Ltime),
		logLevel: logLevel,
		name:     name,
	}, nil
}

func (l *CustomLogger) Debug(v ...interface{}) {
	if l.logLevel <= Debug {
		l.logger.SetPrefix(l.name + DebugPrefix)
		l.logger.Println(v...)
	}
}

func (l *CustomLogger) Info(v ...interface{}) {
	if l.logLevel <= Info {
		l.logger.SetPrefix(l.name + InfoPrefix)
		l.logger.Println(v...)
	}
}

func (l *CustomLogger) Warn(v ...interface{}) {
	if l.logLevel <= Warn {
		l.logger.SetPrefix(l.name + WarnPrefix)
		l.logger.Println(v...)
	}
}

func (l *CustomLogger) Error(format string, v ...interface{}) {
	if l.logLevel <= Error {
		l.logger.SetPrefix(l.name + ErrorPrefix)
		l.logger.Printf(format, v...)
	}
}

// Fatalf logs a formatted error message and then exits the program.
func (l *CustomLogger) Fatalf(format string, v ...interface{}) {

	message := fmt.Sprintf(format, v...)
	fmt.Fprintln(os.Stderr, message)

	os.Exit(1)
}

// Ensure CustomLogger implements Logger
var _ Logger = (*CustomLogger)(nil)

/*
	Note:

	Add write to log file in fatal
	Add logger file rotation
*/
