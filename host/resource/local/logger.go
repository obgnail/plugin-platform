package local

import (
	"github.com/natefinch/lumberjack"
	common "github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/host/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"time"
)

var _ common.PluginLogger = (*Log)(nil)

var Logger *Log

type Log struct {
	Logger *zap.SugaredLogger
}

func (log *Log) Info(format string) {
	log.Logger.Infof(format)
}

func (log *Log) Debug(format string) {
	log.Logger.Debugf(format)
}

func (log *Log) Warn(format string) {
	log.Logger.Warnf(format)
}

func (log *Log) Error(format string) {
	log.Logger.Errorf(format)
}

func (log *Log) ErrorTrace(format string) {
	log.Logger.Error(format)
}

func (log *Log) Fatal(format string) {
	log.Logger.Fatalf(format)
}

func init() {
	logPath := filepath.Join(config.StringOrPanic("runtime_log_dir"), time.Now().Format("2006-01-02"))
	hook := lumberjack.Logger{
		Filename: logPath, // 日志文件路径
		MaxSize:  128,     // 每个日志文件保存的大小 单位:M
		//MaxAge:     30,      // 文件最多保存多少天
		//MaxBackups: 30,      // 日志文件最多保存多少个备份
		Compress: false, // 是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	// 如果是开发环境，同时在控制台上也输出
	//if debug {
	//	writes = append(writes, zapcore.AddSync(os.Stdout))
	//}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 向上跳一级打印出调用打印日志的地方
	skip := zap.AddCallerSkip(1)
	// 输出调用堆栈
	//stack := zap.AddStacktrace(zapcore.DebugLevel)

	// 构造日志
	ZapLogger := zap.New(core, caller, skip, development)
	logger := ZapLogger.Sugar()
	//ZapLogger.Info("log 初始化成功")

	Logger = &Log{Logger: logger}
}
