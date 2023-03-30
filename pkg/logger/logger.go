/**
 * @Author Nil
 * @Description pkg/logger/logger.go
 * @Date 2023/3/28 14:25
 **/

package logger

import (
	"fmt"
	"log"
	"runtime"

	"github.com/ha5ky/hu5ky-bot/pkg/config"
	"github.com/ha5ky/hu5ky-bot/pkg/util"
	"gorm.io/gorm/logger"
)

// constant map
var (
	LogDebug   = 1 << 0
	LogInfo    = 1 << 1
	LogWarning = 1 << 2
	LogError   = 1 << 3
	LogFatal   = 1 << 4

	LogCode2Des = map[int]string{
		LogDebug:   "[DEBUG] ",
		LogInfo:    "[INFO] ",
		LogWarning: "[WARNING] ",
		LogError:   "[ERROR] ",
		LogFatal:   "[FATAL] ",
	}

	LogDes2Code = map[string]int{
		"debug":   LogDebug,
		"info":    LogInfo,
		"warning": LogWarning,
		"error":   LogError,
		"fatal":   LogFatal,
	}

	MySQLLogLevel = map[string]logger.LogLevel{
		"silent": logger.Silent,
		"error":  logger.Error,
		"warn":   logger.Warn,
		"info":   logger.Info,
	}
)

func write2File(msg any, infoLevel int) {
	// TODO write log file
	if infoLevel >= LogDes2Code[config.SysCache.ServerConfig.LogLevel] {
		log.SetPrefix(LogCode2Des[infoLevel])
		pc, callerFile, line, _ := runtime.Caller(2)
		funcName := runtime.FuncForPC(pc).Name()
		log.Println(fmt.Sprintf("%s %s:%d %s", msg, callerFile, line, funcName))
	}
}

func Error(msgs ...interface{}) {
	msg := util.MergeAny(msgs)
	write2File(msg, LogError)
}

func Errorf(format string, msgs ...interface{}) {
	write2File(fmt.Sprintf(format, msgs...), LogError)
}

func Warning(msgs ...interface{}) {
	msg := util.MergeAny(msgs)
	write2File(msg, LogWarning)
}

func Warningf(format string, msgs ...interface{}) {
	write2File(fmt.Sprintf(format, msgs...), LogWarning)
}

func Info(msgs ...interface{}) {
	msg := util.MergeAny(msgs)
	write2File(msg, LogInfo)
}

func Infof(format string, msgs ...interface{}) {
	write2File(fmt.Sprintf(format, msgs...), LogInfo)
}

func Debug(msgs ...interface{}) {
	msg := util.MergeAny(msgs)
	write2File(msg, LogDebug)
}

func Debugf(format string, msgs ...interface{}) {
	write2File(fmt.Sprintf(format, msgs...), LogDebug)
}

func Fatal(msgs ...interface{}) {
	msg := util.MergeAny(msgs)
	write2File(msg, LogFatal)
}

func Fatalf(format string, msgs ...interface{}) {
	write2File(fmt.Sprintf(format, msgs...), LogFatal)
}
