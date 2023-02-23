package log

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/errors"
	"io"
	"os"
	"runtime"
	"strings"
)

const (
	TraceLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Level int8

func (l Level) CapitalString() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

type Encoder interface {
	setWriter(writer io.Writer) func() error
	write(level Level, path string, trace string, msg string, kvs []interface{})
}

type ErrorHooker interface {
	logError(level Level, path string, err error)
	log(level Level, path string, format string, args ...interface{})
}

type ILogger interface {
	Trace(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	TracePath(path string, format string, args ...interface{})
	InfoPath(path string, format string, args ...interface{})
	WarnPath(path string, format string, args ...interface{})
	ErrorPath(path string, format string, args ...interface{})
	WarnDetailsPath(path string, err error)
	ErrorDetailsPath(path string, err error)
	WithValues(keysAndValues ...interface{}) ILogger
}

type Logger struct {
	sep        string
	pathPrefix string
	position   int
	level      Level
	keyValues  []interface{}
	encoder    Encoder
	errHooker  ErrorHooker
}

type Config struct {
	Sep         string
	PathPrefix  string
	Encoder     Encoder
	ErrorHooker ErrorHooker
	Target      io.Writer
	Level       Level
}

func NewLogger(config Config) (*Logger, func() error, error) {
	l := new(Logger)
	l.position = -1
	l.level = config.Level
	l.sep = config.Sep
	l.pathPrefix = config.PathPrefix
	var writer io.Writer
	if config.Target != nil {
		writer = config.Target
	} else {
		writer = os.Stdout
	}
	if config.Encoder == nil {
		l.encoder = &PlainEncoder{}
	} else {
		l.encoder = config.Encoder
	}
	flush := l.encoder.setWriter(writer)
	if config.ErrorHooker != nil {
		l.errHooker = config.ErrorHooker
	}
	return l, flush, nil
}

func (l *Logger) WithValues(kvs ...interface{}) ILogger {
	return &Logger{
		sep:        l.sep,
		pathPrefix: l.pathPrefix,
		position:   l.position,
		level:      l.level,
		encoder:    l.encoder,
		errHooker:  l.errHooker,
		keyValues:  append(l.keyValues, kvs...),
	}
}

func (l *Logger) Trace(format string, args ...interface{}) {
	l.TracePath(l.GetPath(), format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.InfoPath(l.GetPath(), format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.WarnPath(l.GetPath(), format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.ErrorPath(l.GetPath(), format, args...)
}

func (l *Logger) TracePath(path string, format string, args ...interface{}) {
	l.log(TraceLevel, path, "", format, args...)
}

func (l *Logger) InfoPath(path string, format string, args ...interface{}) {
	l.log(InfoLevel, path, "", format, args...)
}

func (l *Logger) WarnPath(path string, format string, args ...interface{}) {
	l.log(WarnLevel, path, "", format, args...)
	if l.errHooker != nil {
		l.errHooker.log(WarnLevel, path, format, args...)
	}
}

func (l *Logger) ErrorPath(path string, format string, args ...interface{}) {
	l.log(ErrorLevel, path, "", format, args...)
	if l.errHooker != nil {
		l.errHooker.log(ErrorLevel, path, format, args...)
	}
}

func (l *Logger) WarnDetails(err error) {
	l.WarnDetailsPath(l.GetPath(), err)
}

func (l *Logger) WarnDetailsPath(path string, err error) {
	l.detailsPath(WarnLevel, path, err)
}
func (l *Logger) ErrorDetails(err error) {
	l.ErrorDetailsPath(l.GetPath(), err)
}

func (l *Logger) ErrorDetailsPath(path string, err error) {
	l.detailsPath(ErrorLevel, path, err)
}

func (l *Logger) detailsPath(level Level, path string, err error) {
	if err == nil {
		return
	}
	var stacks []string
	if e, ok := err.(*errors.Err); ok {
		stacks = e.StackTrace()
		var errlog strings.Builder
		for i := len(stacks) - 1; i >= 0; i-- {
			stacks[i] = l.trimPathPrefix(stacks[i])
			errlog.WriteString(stacks[i])
			if i > 0 {
				errlog.WriteByte('\n')
			}
		}
		l.log(level, "", errlog.String(), err.Error())
	} else {
		l.log(level, path, "", err.Error())
	}
	if l.errHooker != nil {
		l.errHooker.logError(level, path, err)
	}
}

func (l *Logger) trimPathPrefix(s string) string {
	if !strings.Contains(s, "panic") {
		pathprefix := strings.TrimPrefix(l.pathPrefix, "/")
		index := strings.Index(s, pathprefix)
		if index >= 0 {
			index += len(pathprefix)
			if s[index] == '/' {
				index += 1
			}
			return s[index:]
		} else {
			return s
		}
	} else {
		return s
	}
}

func (l *Logger) log(level Level, path string, trace string, format string, args ...interface{}) {
	if l.level > level {
		return
	}
	l.encoder.write(level, path, trace, fmt.Sprintf(format, args...), l.keyValues)
}

func (l *Logger) GetPath() string {
	if _, file, line, ok := runtime.Caller(2); ok {
		if l.position < 0 {
			l.position = strings.LastIndex(file, l.sep)
		}
		if l.position > 0 && l.position < len(file) {
			file = file[l.position:]
		}
		file = strings.TrimPrefix(file, l.pathPrefix)
		return fmt.Sprintf("%s:%d", file, line)
	}
	return ""
}
