package logger

import (
	context "context"
	"log/slog"
	"testing"
)

func TestLoggerLevel(t *testing.T) {
	ctx := context.Background()
	Logger.Debug("LevelDebug")
	Logger.Info("LevelInfo")
	Logger.Warn("LevelWarn")
	Logger.Error("LevelError")
	Record(ctx, slog.LevelDebug, "LevelDebug with context")
	Record(ctx, slog.LevelInfo, "LevelInfo with context")
	Record(ctx, slog.LevelWarn, "LevelWarn with context")
	Record(ctx, slog.LevelError, "LevelError with context")
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func TestLoggerAttr(t *testing.T) {
	ctx := context.Background()
	product := Product{ID: 1, Name: "商品A", Price: 99.9}
	product1 := Product{ID: 1, Name: "商品A", Price: 99.9}
	RecordAttrs(ctx, slog.LevelInfo, "商品信息", slog.Any("product", product), slog.Any("product1", product1))
}

// 测试带traceId和spanId的日志
func TestLoggerTrace(t *testing.T) {
	// 生成traceId和spanId
	traceId := GenerateTraceId()
	spanId := GenerateSpanId("")

	// 创建带有traceId和spanId的context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceId", traceId)
	ctx = context.WithValue(ctx, "spanId", spanId)

	// 记录多条日志，验证分布式追踪
	RecordAttrs(ctx, slog.LevelInfo, "请求开始",
		slog.String("method", "GET"),
		slog.String("path", "/api/users"))

	RecordAttrs(ctx, slog.LevelDebug, "查询数据库",
		slog.String("table", "users"),
		slog.Int("query_time_ms", 15))

	RecordAttrs(ctx, slog.LevelInfo, "请求完成",
		slog.Int("status_code", 200),
		slog.Int64("duration_ms", 45))

	// 测试嵌套调用的spanId生成
	childSpanId := GenerateSpanId(spanId)
	ctxWithChildSpan := context.WithValue(ctx, "spanId", childSpanId)
	RecordAttrs(ctxWithChildSpan, slog.LevelInfo, "子服务调用",
		slog.String("service", "user-service"),
		slog.String("span_id", childSpanId))

	// 测试第二层嵌套
	grandchildSpanId := GenerateSpanId(childSpanId)
	ctxWithGrandchildSpan := context.WithValue(ctx, "spanId", grandchildSpanId)
	RecordAttrs(ctxWithGrandchildSpan, slog.LevelInfo, "深层调用",
		slog.String("operation", "cache_lookup"),
		slog.String("span_id", grandchildSpanId))
}

// 测试多个属性
func TestLoggerAttrs(t *testing.T) {
	ctx := context.Background()
	RecordAttrs(ctx, slog.LevelInfo, "系统状态",
		slog.String("host", "server-01"),
		slog.Int("cpu_usage", 75),
		slog.Float64("memory_usage", 82.5),
		slog.String("status", "running"),
		slog.Group("database",
			slog.String("host", "db.example.com"),
			slog.Int("port", 5432),
			slog.String("name", "mydb"),
		),
	)
}
