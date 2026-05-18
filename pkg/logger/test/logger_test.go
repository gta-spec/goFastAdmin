package test

import (
	"gota/pkg/logger"
	"gota/pkg/logger/rotate"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestNewRollerTime(t *testing.T) {
	roller, err := rotate.NewRoller("runtime/log/%Y-%m/%d.log", int64(512), &rotate.Options{
		MaxBackups: 16,
		MaxAge:     24 * time.Hour,
		LocalTime:  true,
		Compress:   false,
	})
	if err != nil {
		t.Fatalf("Failed to create roller: %v", err)
	}
	defer roller.Close()

	// 创建 logger
	testLogger := logger.New(
		logger.WithLevel(slog.LevelDebug),
		logger.WithCallerSkip(1),
		logger.WithWriter(roller),
	)
	// 写入足够多的数据以触发多次轮转
	// 每条日志大约 100-150 字节，写入 50 条约 5000-7500 字节
	// 应该能触发 10-15 次轮转（512字节限制）
	writeCount := 50
	for i := 0; i < writeCount; i++ {
		testLogger.Info("测试日志轮转功能", slog.Int("iteration", i), slog.String("data", strings.Repeat("x", 50)))
	}
	// 等待异步的 mill 操作完成
	time.Sleep(100 * time.Millisecond)
}
