package logger

import (
	"io"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

var (
	// Log 全局日志实例
	Log  *logrus.Logger
	once sync.Once
)

// Config 包含日志配置信息
type Config struct {
	Level      string // 日志级别: debug, info, warn, error, fatal, panic
	Filename   string // 日志文件路径
	MaxSize    int    // 单个日志文件最大尺寸(MB)
	MaxBackups int    // 保留的旧日志文件最大数量
	MaxAge     int    // 旧日志文件保留的最大天数
	Compress   bool   // 是否压缩旧日志文件
	Console    bool   // 是否输出到控制台
	Format     string // 输出格式: json, text
}

//	"level": "debug",
//
// "filename": "logs/app.log",
// "max_size": 100,
// "max_backups": 3,
// "max_age": 28,
// "compress": true,
// "console": true,
// "format": "json"
var DefaultConfig = &Config{
	Level:      "debug",
	Filename:   "logs/app.log",
	MaxSize:    10,
	MaxBackups: 5,
	MaxAge:     30,
	Compress:   false,
	Console:    true,
	Format:     "text",
}

// InitLogger 初始化日志
func InitLogger(config *Config) {
	once.Do(func() {
		Log = logrus.New()
		// 设置日志级别
		level, err := logrus.ParseLevel(config.Level)
		if err != nil {
			level = logrus.InfoLevel
		}
		Log.SetLevel(level)

		// 设置格式化器
		if config.Format == "json" {
			Log.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: time.RFC3339,
				PrettyPrint:     false,
			})
		} else {
			Log.SetFormatter(&logrus.TextFormatter{
				TimestampFormat: time.RFC3339,
				FullTimestamp:   true,
			})
		}
		// 添加调用者信息
		Log.SetReportCaller(true)
		// 设置输出
		var writers []io.Writer
		// 控制台输出
		if config.Console {
			writers = append(writers, os.Stdout)
		}
		// 文件输出
		if config.Filename != "" {
			// 确保日志目录存在
			logDir := path.Dir(config.Filename)
			if err := os.MkdirAll(logDir, 0755); err != nil {
				Log.Errorf("无法创建日志目录: %v", err)
			}

			// 配置日志轮转
			rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
				Filename:   config.Filename,
				MaxSize:    config.MaxSize,
				MaxBackups: config.MaxBackups,
				MaxAge:     config.MaxAge,
				Level:      level,
				Formatter: &logrus.JSONFormatter{
					TimestampFormat: time.RFC3339,
				},
			})

			if err != nil {
				Log.Errorf("无法初始化文件日志: %v", err)
			} else {
				Log.AddHook(rotateFileHook)
			}
		}

		// 如果有多个输出，设置多写入器
		if len(writers) > 0 {
			if len(writers) == 1 {
				Log.SetOutput(writers[0])
			} else {
				mw := io.MultiWriter(writers...)
				Log.SetOutput(mw)
			}
		}

		Log.Info("Logrus 日志系统初始化成功")
	})
}

// 提供便捷方法
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}
