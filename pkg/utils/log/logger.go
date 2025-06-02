// 文件路径：douyin/pkg/utils/log/logger.go
package log

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

// LogrusObj 全局日志对象
var LogrusObj *logrus.Logger

// setOutputFile 设置日志文件输出，返回日志文件的 os.File 对象
func setOutputFile() (*os.File, error) {
	now := time.Now()
	var logFilePath string
	// 获取当前工作目录，构造日志目录路径
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	// 检查日志目录是否存在，不存在则创建
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		if err := os.MkdirAll(logFilePath, 0777); err != nil {
			fmt.Println("创建日志目录失败：" + err.Error())
			return nil, err
		}
	}
	// 构造日志文件名称，如 "2025-02-21.log"
	logFileName := now.Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	// 检查日志文件是否存在，不存在则创建
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println("创建日志文件失败：" + err.Error())
			return nil, err
		}
	}
	// 打开日志文件，追加写入
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("打开日志文件失败：" + err.Error())
		return nil, err
	}
	return src, nil
}

// InitLog 初始化日志
func InitLog() {
	if LogrusObj == nil {
		LogrusObj = logrus.New()
		LogrusObj.SetLevel(logrus.DebugLevel)         // 设置日志级别
		LogrusObj.SetFormatter(&logrus.TextFormatter{ // 设置文本格式
			FullTimestamp: true,
		})
		// 添加文件输出逻辑（兼容第一版）
		if file, err := setOutputFile(); err == nil {
			LogrusObj.Out = file
		} else {
			// 如果文件设置失败，则默认输出到标准输出
			fmt.Println("使用默认标准输出，文件输出初始化失败：", err)
		}
	}
}

// Info 输出信息级别的日志
func Info(message string) {
	if LogrusObj == nil {
		InitLog()
	}
	LogrusObj.Info(message)
}

// Errorf 输出格式化的错误日志
func Errorf(format string, args ...interface{}) {
	if LogrusObj == nil {
		InitLog()
	}
	LogrusObj.Errorf(format, args...)
}

// Infof 输出格式化的信息日志
func Infof(format string, args ...interface{}) {
	if LogrusObj == nil {
		InitLog()
	}
	LogrusObj.Infof(format, args...)
}
