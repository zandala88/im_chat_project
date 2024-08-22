package public

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

func init() {
	lumberjacklogger := &lumberjack.Logger{
		Filename:   "./runlog/log-rotate.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	defer lumberjacklogger.Close()

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder // 设置时间格式
	fileEncoder := zapcore.NewConsoleEncoder(config)
	core := zapcore.NewCore(
		fileEncoder,                       //编码设置
		zapcore.AddSync(lumberjacklogger), //输出到文件
		zap.DebugLevel,                    //日志等级
	)

	log := zap.New(core, zap.AddCaller())
	defer log.Sync()
	zap.ReplaceGlobals(log)
}

type CustomLogger struct {
	logLevel                  logger.LogLevel
	ignoreRecordNotFoundError bool
	slowThreshold             time.Duration
}

func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *CustomLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.logLevel >= logger.Info {
		zap.S().Debugf(msg, args...)
	}
}

func (l *CustomLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.logLevel >= logger.Warn {
		zap.S().Infof(msg, args...)
	}
}

func (l *CustomLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.logLevel >= logger.Error {
		zap.S().Errorf(msg, args...)
	}
}

func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if l.ignoreRecordNotFoundError && errors.Is(err, gorm.ErrRecordNotFound) {
		// Ignore record not found errors
		return
	}

	switch {
	case err != nil && l.logLevel >= logger.Error:
		zap.S().Errorf("err: %v | elapsed: %v | sql: %s | rows: %d", err, elapsed, sql, rows)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
		zap.S().Infof("elapsed: %v | slow sql: %s | rows: %d", elapsed, sql, rows)
	case l.logLevel >= logger.Info:
		zap.S().Debugf("elapsed: %v | sql: %s | rows: %d", elapsed, sql, rows)
	}
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				zap.L().Error("[Recovery from panic]",
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
