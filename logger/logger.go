package logger

import (
	"io"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var errorLogger *zap.SugaredLogger

func init() {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:       "created_at",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "from",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		// EncodeLevel: zapcore.CapitalLevelEncoder,
		// EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	// 实现两个判断日志等级的interface
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	debugWriter := getWriter("./logs/debug.log")
	infoWriter := getWriter("./logs/info.log")
	warnWriter := getWriter("./logs/warn.log")
	errorWriter := getWriter("./logs/error.log")
	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	errorLogger = log.Sugar()
}
func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1) + "-%Y%m%d%H.log", // 没有使用go风格反人类的format格式
	//rotatelogs.WithLinkName(filename),
	//rotatelogs.WithMaxAge(time.Hour*24*7),
	//rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return hook
}
func Debug(args ...interface{}) {
	errorLogger.Debug(args...)
}
func Debugf(template string, args ...interface{}) {
	errorLogger.Debugf(template, args...)
}
func Info(args ...interface{}) {
	errorLogger.Info(args...)
}
func Infof(template string, args ...interface{}) {
	errorLogger.Infof(template, args...)
}
func Warn(args ...interface{}) {
	errorLogger.Warn(args...)
}
func Warnf(template string, args ...interface{}) {
	errorLogger.Warnf(template, args...)
}
func Error(args ...interface{}) {
	errorLogger.Error(args...)
}
func Errorf(template string, args ...interface{}) {
	errorLogger.Errorf(template, args...)
}
func DPanic(args ...interface{}) {
	errorLogger.DPanic(args...)
}
func DPanicf(template string, args ...interface{}) {
	errorLogger.DPanicf(template, args...)
}
func Panic(args ...interface{}) {
	errorLogger.Panic(args...)
}
func Panicf(template string, args ...interface{}) {
	errorLogger.Panicf(template, args...)
}
func Fatal(args ...interface{}) {
	errorLogger.Fatal(args...)
}
func Fatalf(template string, args ...interface{}) {
	errorLogger.Fatalf(template, args...)
}
