// 文件路径：douyin/pkg/utils/log/logger.go
package log

import (
	"fmt"
	"io" // Needed for io.Discard
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin" // Used for gin.DefaultErrorWriter if needed as fallback
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// LogrusObj 全局日志对象
var LogrusObj *logrus.Logger // Or simply: var Log *logrus.Logger

// InitLogger 初始化日志
func InitLogger() {
	if LogrusObj != nil {
		// Already initialized
		return
	}

	LogrusObj = logrus.New()
	LogrusObj.SetLevel(logrus.DebugLevel) // Set appropriate log level from config eventually

	// Set JSON Formatter
	LogrusObj.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339, // e.g., "2006-01-02T15:04:05Z07:00"
		// FieldMap: logrus.FieldMap{ // Optional: customize field names
		// 	logrus.FieldKeyTime:  "timestamp",
		// 	logrus.FieldKeyLevel: "level",
		// 	logrus.FieldKeyMsg:   "message",
		// 	logrus.FieldKeyFunc:  "caller", // Requires LogrusObj.SetReportCaller(true)
		// },
	})
	// LogrusObj.SetReportCaller(true) // Uncomment if you want to log filename and line number

	// Configure file-rotatelogs
	logFilePath := "logs" // Base directory for logs
	logFileName := "app.log"
	logFile := path.Join(logFilePath, logFileName)

	// Ensure logs directory exists
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		if errMkdir := os.MkdirAll(logFilePath, 0755); errMkdir != nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "创建日志目录失败：%s\n", errMkdir.Error())
			// Fallback to stdout if directory creation fails
			LogrusObj.SetOutput(os.Stdout)
			return
		}
	}

	writer, err := rotatelogs.New(
		logFile+".%Y%m%d",                 // Rotated file name pattern
		rotatelogs.WithLinkName(logFile),          // Link to current log file
		rotatelogs.WithMaxAge(7*24*time.Hour),     // Max age of 7 days
		rotatelogs.WithRotationTime(24*time.Hour), // Rotate daily
		// rotatelogs.WithRotationSize(100*1024*1024), // Optional: Rotate by size
	)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "初始化 rotatelogs 失败：%s\n", err.Error())
		// Fallback to stdout if rotatelogs setup fails
		LogrusObj.SetOutput(os.Stdout)
		return
	}

	// Create a new LFS hook
	hook := lfshook.NewHook(
		lfshook.WriterMap{ // Log all levels to the rotatelogs writer
			logrus.TraceLevel: writer,
			logrus.DebugLevel: writer,
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		},
		&logrus.JSONFormatter{TimestampFormat: time.RFC3339}, // Ensure hook uses the same formatter
	)
	LogrusObj.AddHook(hook)

	// Discard default Logrus output (os.Stderr), rely entirely on the hook for file logging.
	// If console output is also desired for development, you can add another hook
	// or set LogrusObj.SetOutput(io.MultiWriter(os.Stdout, writer_from_hook_or_another_rotatelog))
	// but that makes the hook for specific levels to 'writer' a bit redundant if 'writer' is also default.
	// For dedicated file logging via hook, and potentially a separate console hook if needed:
	LogrusObj.SetOutput(io.Discard)

	// Use fmt.Println for initial messages that should go to console regardless of log setup
	fmt.Println("Logrus 日志系统初始化成功，使用 JSON 格式和每日轮转。日志文件位于:", logFile)
}

// Helper functions - should not call InitLogger themselves.
// Assume LogrusObj is initialized at startup.

// Info 输出信息级别的日志
func Info(message string) {
	if LogrusObj == nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "警告: LogrusObj 未初始化，调用 Info 失败。消息: %s\n", message)
		return
	}
	LogrusObj.Info(message)
}

// Errorf 输出格式化的错误日志
func Errorf(format string, args ...interface{}) {
	if LogrusObj == nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "警告: LogrusObj 未初始化，调用 Errorf 失败。格式: %s, 参数: %v\n", format, args)
		return
	}
	LogrusObj.Errorf(format, args...)
}

// Infof 输出格式化的信息日志
func Infof(format string, args ...interface{}) {
	if LogrusObj == nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "警告: LogrusObj 未初始化，调用 Infof 失败。格式: %s, 参数: %v\n", format, args)
		return
	}
	LogrusObj.Infof(format, args...)
}

// Warnf outputs a formatted warning log
func Warnf(format string, args ...interface{}) {
	if LogrusObj == nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "警告: LogrusObj 未初始化，调用 Warnf 失败。格式: %s, 参数: %v\n", format, args)
		return
	}
	LogrusObj.Warnf(format, args...)
}

// Debugf outputs a formatted debug log
func Debugf(format string, args ...interface{}) {
	if LogrusObj == nil {
		// This might be too noisy for production if logger isn't init, but useful for dev.
		// Consider removing if it becomes spammy.
		fmt.Fprintf(gin.DefaultErrorWriter, "警告: LogrusObj 未初始化，调用 Debugf 失败。格式: %s, 参数: %v\n", format, args)
		return
	}
	LogrusObj.Debugf(format, args...)
}

// Panic logs a message at panic level and then panics.
func Panic(message string) {
	if LogrusObj == nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "警告: LogrusObj 未初始化，调用 Panic 失败。消息: %s\n", message)
		panic(message) // Still panic
	}
	LogrusObj.Panic(message)
}
