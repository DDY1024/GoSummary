package main

import (
	"io"

	// rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// func main() {
// 	// time="2023-06-05T17:38:54+08:00" level=info msg="A walrus appears" animal=walrus
// 	log.WithFields(log.Fields{
// 		"animal": "walrus",
// 	}).Info("A walrus appears")
// }

var (
	infoWriter  io.Writer
	errorWriter io.Writer
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	// log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	// log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	// log.SetLevel(log.WarnLevel)

	// out.log
	// infoWriter, _ = rotatelogs.New(
	// 	"out.log"+".%Y%m%d%H",
	// 	rotatelogs.WithLinkName("out.log"),
	// 	rotatelogs.WithRotationTime(time.Hour),
	// 	rotatelogs.WithMaxAge(7*24*time.Hour),
	// )

	// // out.error.log
	// errorWriter, _ = rotatelogs.New(
	// 	"out.error.log"+".%Y%m%d%H",
	// 	rotatelogs.WithLinkName("out.error.log"),
	// 	rotatelogs.WithRotationTime(time.Hour),
	// 	rotatelogs.WithMaxAge(7*24*time.Hour),
	// )

	// log.SetOutput(writer)

	// log.Hook.Add()
}

func testLog() {

	// lg := log.New()
	// lg.Hooks.Add(CustomHook{})

	logrus.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	logrus.WithFields(log.Fields{
		"omg":    true,
		"number": 122,
	}).Error("The group's number increased tremendously!")

	// exit status 1
	// 严重错误，需要停止程序
	// log.WithFields(log.Fields{
	// 	"omg":    true,
	// 	"number": 100,
	// }).Fatal("The ice breaks!")

	// log.WithFields(log.Fields{
	// 	"omg":    true,
	// 	"number": 100,
	// }).Panic("The ice breaks!")

	// nlog := log.New()
	// nlog.WithFields()

	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	contextLogger := logrus.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
		"xxx":    "xxx",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Error("Me too")

	// lg := logrus.New()
	// lg.Hooks.Add()

	// lg := logrus.New()
	// lg.Hooks.Add()
}

/*
// Levels
// Fire
type Hook interface {
	Levels() []Level
	Fire(*Entry) error
}
*/

type CustomHook struct{}

func (h CustomHook) Levels() []log.Level {
	return log.AllLevels
}

func (h CustomHook) Fire(entry *log.Entry) error {
	// bs, err := entry.Bytes()
	// if err != nil {
	// 	return err
	// }

	// 性能存在问题: 存在加锁操作，性能损耗
	// 本身内部存在加锁操作，应该还好
	// 如果对性能有要求，考虑 logger 分片，但是并发写文件又会存在问题
	switch entry.Level {
	case log.ErrorLevel:
		// entry.Info()
		entry.Logger.SetOutput(errorWriter)
		// errorWriter.Write(bs)
	default:
		entry.Logger.SetOutput(infoWriter)
	}
	return nil
}

// 特殊符号分隔
// x,y,z
// x:y:z
