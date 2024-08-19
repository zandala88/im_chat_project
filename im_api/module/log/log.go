package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	lumberjacklogger := &lumberjack.Logger{
		Filename:   "./runlog/log-rotate-test.json",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	defer lumberjacklogger.Close()

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder // 设置时间格式
	fileEncoder := zapcore.NewJSONEncoder(config)
	core := zapcore.NewCore(
		fileEncoder,                       //编码设置
		zapcore.AddSync(lumberjacklogger), //输出到文件
		zap.DebugLevel,                    //日志等级
	)

	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
}
