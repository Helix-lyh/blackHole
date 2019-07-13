package Logger

import (
	"fmt"
	"github.com/op/go-logging"
	"io"
	"os"
)


// Password只是实现编校器接口的示例类型。任何
// 记录此日志时，将调用 Redacted()函数。
type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

func GetLogger(name string) *logging.Logger {
	var logger = logging.MustGetLogger(name)
	//示例格式字符串。除了消息之外，所有内容都有自定义颜色
	//取决于日志级别。许多字段都有自定义输出
	//格式化也一样。时间返回到毫秒。
	var format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	// 为os.Stderr创建两个后端.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	_, err := os.Stat("glog.log")
	var logFile  io.Writer
	if err != nil{
		logFile, err = os.Create("glog.log")
		if err != nil{
			fmt.Printf("日志文件创建异常! Error=%v\n", err)
		}
	}else {
		logFile, err = os.OpenFile("glog.log", os.O_WRONLY, 0644)
		if err != nil{
			fmt.Printf("日志文件打开异常! Error=%v\n", err)
		}
	}
	backend3 := logging.NewLogBackend(logFile,"", 0)

	//写入backend2的消息，添加一些额外的内容 包括使用的日志级别和名称
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend3Formatter := logging.NewBackendFormatter(backend3, format)

	// 错误和更严重的消息才发送到backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// 设置要使用的后端
	logging.SetBackend(backend1Leveled, backend2Formatter, backend3Formatter)

	//logger.Debugf("debug %s", Password("secret"))
	//logger.Info("info")
	//logger.Notice("notice")
	//logger.Warning("warning")
	//logger.Error("err")
	//logger.Critical("crit")
	logger.Info("logger初始化完成!")
	return logger
}