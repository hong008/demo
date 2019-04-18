package logrus

import (
	"demo/5.go_log/base"
	"demo/5.go_log/field"
	"github.com/Sirupsen/logrus"
	"io"
	"os"
)

type loggerLogrus struct {
	level           base.LogLevel
	format          base.LogFormat
	optWithLocation base.OptWithLocation
	inner           *logrus.Entry
}

func NewLogger() base.MyLogger {
	return NewLoggerBy(base.LEVEL_INFO, base.FORMAT_TEXT, os.Stdout, nil)
}

func NewLoggerBy(level base.LogLevel, format base.LogFormat, writer io.Writer, options []base.Option) base.MyLogger {
	var logrusLevel logrus.Level
	switch level {
	case base.LEVEL_DEBUG:
		logrusLevel = logrus.DebugLevel
	case base.LEVEL_WARN:
		logrusLevel = logrus.WarnLevel
	case base.LEVEL_ERROR:
		logrusLevel = logrus.ErrorLevel
	case base.LEVEL_FATAL:
		logrusLevel = logrus.FatalLevel
	case base.LEVEL_PANIC:
		logrusLevel = logrus.PanicLevel
	default:
		level = base.LEVEL_INFO
		logrusLevel = logrus.InfoLevel
	}

	var optWithLocation base.OptWithLocation
	if options != nil {
		for _, opt := range options {
			if opt.Name() == "with location" {
				optWithLocation, _ = opt.(base.OptWithLocation)
			}
		}
	}
	return &loggerLogrus{
		level:           level,
		format:          format,
		optWithLocation: optWithLocation,
		inner:           initInnerLogger(logrusLevel, format, writer),
	}
}

func (logger *loggerLogrus) getInner() *logrus.Entry {
	return logger.inner
}

func WithLocation(entry *logrus.Entry) *logrus.Entry {
	funcPath, fileName, line := base.GetInvokeerLocation(4)
	return entry.WithField(
		"location", map[string]interface{}{
			"func_path": funcPath,
			"file_name": fileName,
			"line":      line,
		})
}

func initInnerLogger(level logrus.Level, format base.LogFormat, writer io.Writer) *logrus.Entry {
	innerLogger := logrus.New()

	switch format {
	case base.FORMAT_JSON:
		innerLogger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: base.TIMESTAMP_FORMAT,
		}
	default:
		innerLogger.Formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: base.TIMESTAMP_FORMAT,
			DisableSorting:  false,
		}
	}
	innerLogger.Level = level
	innerLogger.Out = writer
	return logrus.NewEntry(innerLogger)
}

func (logger *loggerLogrus) Name() string {
	return "logrus"
}

func (logger *loggerLogrus) Level() base.LogLevel {
	return logger.level
}

func (logger *loggerLogrus) Format() base.LogFormat {
	return logger.format
}

func (logger *loggerLogrus) Options() []base.Option {
	return []base.Option{logger.optWithLocation}
}

func (logger *loggerLogrus) Debug(v ...interface{}) {
	logger.getInner().Debug(v...)
}

func (logger *loggerLogrus) Debugf(format string, v ...interface{}) {
	logger.getInner().Debugf(format, v...)
}

func (logger *loggerLogrus) Debugln(v ...interface{}) {
	logger.getInner().Debugln(v...)
}

func (logger *loggerLogrus) Error(v ...interface{}) {
	logger.getInner().Error(v...)
}

func (logger *loggerLogrus) Errorf(format string, v ...interface{}) {
	logger.getInner().Errorf(format, v...)
}

func (logger *loggerLogrus) Errorln(v ...interface{}) {
	logger.getInner().Errorln(v...)
}

func (logger *loggerLogrus) Fatal(v ...interface{}) {
	logger.getInner().Fatal(v...)
}

func (logger *loggerLogrus) Fatalf(format string, v ...interface{}) {
	logger.getInner().Fatalf(format, v...)
}

func (logger *loggerLogrus) Fatalln(v ...interface{}) {
	logger.getInner().Fatalln(v...)
}

func (logger *loggerLogrus) Info(v ...interface{}) {
	logger.getInner().Info(v...)
}

func (logger *loggerLogrus) Infof(format string, v ...interface{}) {
	logger.getInner().Infof(format, v...)
}

func (logger *loggerLogrus) Infoln(v ...interface{}) {
	logger.getInner().Infoln(v...)
}

func (logger *loggerLogrus) Warn(v ...interface{}) {
	logger.getInner().Warn(v...)
}

func (logger *loggerLogrus) Warnf(format string, v ...interface{}) {
	logger.getInner().Warnf(format, v...)
}

func (logger *loggerLogrus) Warnln(v ...interface{}) {
	logger.getInner().Warnln(v...)
}

func (logger *loggerLogrus) Panic(v ...interface{}) {
	logger.getInner().Panic(v...)
}

func (logger *loggerLogrus) Panicf(format string, v ...interface{}) {
	logger.getInner().Panicf(format, v...)
}

func (logger *loggerLogrus) Panicln(v ...interface{}) {
	logger.getInner().Panicln(v...)
}

func (logger *loggerLogrus) WithFields(fields ...field.Field) base.MyLogger {
	fieldsLen := len(fields)
	if fieldsLen == 0 {
		return logger
	}

	logrusFields := make(map[string]interface{}, fieldsLen)
	for _, curfield := range fields {
		logrusFields[curfield.Name()] = curfield.Value()
	}

	return &loggerLogrus{
		level:           logger.level,
		format:          logger.format,
		optWithLocation: logger.optWithLocation,
		inner:           logger.inner.WithFields(logrusFields),
	}
}
