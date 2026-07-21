package zap

import (
	"fmt"
	"io"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log *zap.Logger
}

type option struct {
	maxAge int
	encode *zapcore.EncoderConfig
}
type Option func(*option)

func WithMaxAge(day int) Option {
	return func(o *option) {
		o.maxAge = day
	}
}
func WithEncoderConfig(encode *zapcore.EncoderConfig) Option {
	return func(o *option) {
		o.encode = encode
	}
}

func NewFsLogger(path string, opts ...Option) log.Logger {

	// "ts", log.DefaultTimestamp,
	// 	"caller", log.DefaultCaller,
	// 	"service.id", id,
	// 	"service.name", Name,
	// 	"service.version", Version,
	// 	"trace_id", log.TraceID(),
	// 	"span_id", log.SpanID(),

	opt := &option{
		maxAge: 7,
		encode: &zapcore.EncoderConfig{
			TimeKey:        "time",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,

			NameKey:       "logger",
			MessageKey:    "msg",
			StacktraceKey: "stack",
			EncodeName:    zapcore.FullNameEncoder,
		},
	}

	for _, o := range opts {
		o(opt)
	}

	// hook, err := rotatelogs.New(
	// 	path,
	// 	rotatelogs.WithMaxAge(time.Hour*24*time.Duration(opt.maxAge)),
	// )

	// if err != nil {
	// 	panic(err)
	// }

	encoder := zapcore.NewJSONEncoder(*opt.encode)

	core := zapcore.NewTee(
		zapcore.NewCore(
			encoder,
			zapcore.AddSync(&lumberjack.Logger{
				Filename: path,       //指定日志存储位置
				MaxAge:   opt.maxAge, //日志文件存储最大天数
				Compress: false,      //是否执行压缩
				// MaxSize:    100,        //日志的最大大小（M）
				// MaxBackups: 50,         //日志的最大保存数量
			}),
			zap.InfoLevel,
		),
	)

	// // 打印
	// core := zapcore.NewCore(
	// 	encoder,
	// 	zapcore.NewMultiWriteSyncer(
	// 		zapcore.AddSync(os.Stdout),
	// 	),
	// 	zap.InfoLevel)

	return &Logger{
		log: zap.New(
			core,
			zap.AddCaller(),
			zap.AddCallerSkip(3),
			zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		),
	}
}

func NewFsLoggerWithDay(path, err_path string, opts ...Option) log.Logger {

	// 基本格式设置
	opt := &option{
		maxAge: 7,
		encode: &zapcore.EncoderConfig{
			TimeKey:        "time",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,

			NameKey:       "logger",
			MessageKey:    "msg",
			StacktraceKey: "stack",
			EncodeName:    zapcore.FullNameEncoder,
		},
	}

	for _, o := range opts {
		o(opt)
	}

	// 判断日志等级
	infoLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv < zapcore.WarnLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.WarnLevel
	})

	infoWriter := getWriter(path, opt.maxAge)
	warnWriter := getWriter(err_path, opt.maxAge)

	encoder := zapcore.NewJSONEncoder(*opt.encode)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
	)

	return &Logger{
		log: zap.New(
			core,
			zap.AddCaller(),
			zap.AddCallerSkip(3),
			zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		),
	}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {

	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug("", data...)
	case log.LevelInfo:
		l.log.Info("", data...)
	case log.LevelWarn:
		l.log.Warn("", data...)
	case log.LevelError:
		l.log.Error("", data...)
	case log.LevelFatal:
		l.log.Fatal("", data...)
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func getWriter(path string, day int) io.Writer {
	w, err := rotatelogs.New(
		path+"_%Y%m%d.log",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(day)),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	return w
}
