package core

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var Log *logrus.Logger


func IsTraceEnabled() bool{
	return Log.GetLevel() <= logrus.TraceLevel
}

func IsDebugEnabled() bool{
	return Log.GetLevel() <= logrus.DebugLevel
}

func IsInfoEnabled() bool{
	return Log.GetLevel() <= logrus.InfoLevel
}

func IsWarnEnabled() bool{
	return Log.GetLevel() <= logrus.WarnLevel
}

func IsErrorEnabled() bool{
	return Log.GetLevel() <= logrus.ErrorLevel
}

func IsFatalEnabled() bool{
	return Log.GetLevel() <= logrus.FatalLevel
}

func IsPanicEnabled() bool{
	return Log.GetLevel() <= logrus.PanicLevel
}

//初始化日志
func InitLog(config *Config){
	var maxAgeDuration time.Duration
	var rotationTimeDuration time.Duration
	maxAge := config.Logger.MaxAge
	if maxAge <= 0 {
		maxAgeDuration = time.Hour * 24 * 60		//60天
	}else{
		maxAgeDuration = time.Minute * maxAge
	}
	rotationTime := config.Logger.RotationTime
	if rotationTime <= 0 {
		rotationTimeDuration = time.Hour * 24	//1天
	}else{
		rotationTimeDuration = time.Minute * rotationTime
	}
	configLocalFilesystemLogger(config.Logger.Level,
		config.Logger.Path,
		config.Logger.Filename,
		maxAgeDuration,
		rotationTimeDuration)
}

// config logrus log to local filesystem, with file rotation
func configLocalFilesystemLogger(level string, logPath string, logFileName string,
				maxAge time.Duration, rotationTime time.Duration) {
	if !IsExists(logPath) {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err!=nil {
			panic(err)
		}
	}
	Log = logrus.New()

	lv, err := logrus.ParseLevel(level)
	if err!=nil {
		Log.Panicf("cannot parse log level: %s, %+v", level, errors.WithStack(err))
	}
	Log.SetLevel(lv)

	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d",
		rotatelogs.WithLinkName(baseLogPath), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		Log.Panicf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.TraceLevel: writer,
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{})
	Log.AddHook(lfHook)
}