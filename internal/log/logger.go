// Package log provides configurable logging. It will detect if the process is
// running in kubernetes by searching for the "KUBERNETES_SERVICE_HOST"
// environment variable. If it is running in kubernetes it will output logs to
// stdout using json. If it is not running in kubernetes it will output logs in
// a standard single line readable format.
//
// Additionally, you can set a LOG_LEVEL environment value to any of the
// following values, to retrieve only log levels from that level and above. The
// default log level is INFO for running in kubernetes and DEBUG when not.
//
// FATAL
// ERROR
// WARN
// INFO
// DEBUG
package log

import (
	"strings"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	configuredLogger *zap.SugaredLogger
	mutex            sync.Mutex
)

func InitializeLogger(isTerminal bool, logLevel string) error {
	mutex.Lock()
	defer mutex.Unlock()

	var config zap.Config
	if isTerminal {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	switch strings.ToLower(logLevel) {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		config.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return errors.Wrap(err, "failed to initialize logger")
	}

	configuredLogger = logger.Sugar()
	return nil
}

func init() {
	InitializeLogger(true, "undefined")
}

// Debug logs a message with some additional context.
func Debug(msg string, keysAndValues ...interface{}) {
	configuredLogger.Debugw(msg, keysAndValues...)
}

// Info logs a message with some additional context.
func Info(msg string, keysAndValues ...interface{}) {
	configuredLogger.Infow(msg, keysAndValues...)
}

// Warn logs a message with some additional context.
func Warn(msg string, keysAndValues ...interface{}) {
	configuredLogger.Warnw(msg, keysAndValues...)
}

// Error logs a message with some additional context.
func Error(msg string, keysAndValues ...interface{}) {
	configuredLogger.Errorw(msg, keysAndValues...)
}

// Fatal logs a message with some additional context, then calls os.Exit.
func Fatal(msg string, keysAndValues ...interface{}) {
	configuredLogger.Fatalw(msg, keysAndValues...)
}
