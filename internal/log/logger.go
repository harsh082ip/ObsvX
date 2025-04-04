package log

import (
	"io"
	"os"
	"path/filepath"

	"github.com/harsh082ip/ObsvX/config"
	"github.com/rs/zerolog"
)

const (
	loggerKey = "log_data"
)

var (
	DefaultLoggerAddr *zerolog.Logger
)

type Logger struct {
	Log           *zerolog.Logger
	DefaultLogger *DefaultLog
}

type DefaultLog struct {
	DeviceID   string `json:"device_id,omitempty"`
	RemoteAddr string `json:"remote_addr,omitempty"`
	Source     string `json:"source,omitempty"`
}

// InitLogger creates a new logger instance with the given source
func InitLogger(source string) *Logger {
	if DefaultLoggerAddr == nil {
		defaultLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		defaultLogger.Warn().Msg("Start writing logs on standard output")
		defaultLogger.Error().Msg("Failed to initialize logger.Err: Default logger address be nil")
		DefaultLoggerAddr = &defaultLogger
	}
	logger := DefaultLoggerAddr.With().Logger()
	defaultLog := &DefaultLog{Source: source}
	return &Logger{Log: &logger, DefaultLogger: defaultLog}
}

// SetupGlobalLogger initializes the global logger with configuration
func SetupGlobalLogger(cfg *config.AppConfig) error {
	// Set global log level
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)

	// Setup log output (file and/or console)
	var writers []io.Writer

	// Configure file output if log file is specified
	if cfg.LogFile != "" {
		// Ensure directory exists
		logDir := filepath.Dir(cfg.LogFile)
		if logDir != "." && logDir != "" {
			if err := os.MkdirAll(logDir, 0755); err != nil {
				return err
			}
		}

		// Open log file
		logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		writers = append(writers, logFile)
	}

	// Add console writer if enabled
	if cfg.LogToConsole {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02T15:04:05Z07:00",
		}
		writers = append(writers, consoleWriter)
	}

	// Create multi-writer if we have multiple outputs
	var output io.Writer
	if len(writers) > 1 {
		output = io.MultiWriter(writers...)
	} else if len(writers) == 1 {
		output = writers[0]
	} else {
		// Default to stdout if no writers configured
		output = os.Stdout
	}

	// Create the global logger
	defaultLogger := zerolog.New(output).With().Timestamp().Logger()
	DefaultLoggerAddr = &defaultLogger

	return nil
}

func (l Logger) LogInfoMessage() *zerolog.Event {
	return l.Log.Info().Interface(loggerKey, l.DefaultLogger)
}

func (l Logger) LogWarnMessage() *zerolog.Event {
	return l.Log.Warn().Interface(loggerKey, l.DefaultLogger)
}

func (l Logger) LogDebugMessage() *zerolog.Event {
	return l.Log.Debug().Interface(loggerKey, l.DefaultLogger)
}

func (l Logger) LogErrorMessage() *zerolog.Event {
	return l.Log.Error().Interface(loggerKey, l.DefaultLogger)
}

func (l Logger) LogFatalMessage() *zerolog.Event {
	return l.Log.Fatal().Interface(loggerKey, l.DefaultLogger)
}

func (l *Logger) Print(val string) {
	l.Log.Info().Interface(loggerKey, l.DefaultLogger).Msg(val)
}

func (l *Logger) Printf(format string, val string) {
	l.Log.Info().Interface(loggerKey, l.DefaultLogger).Msg(val)
}

func (l *Logger) Println(val string) {
	l.Log.Info().Interface(loggerKey, l.DefaultLogger).Msg(val)
}
