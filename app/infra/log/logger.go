package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zapcore.Field)
	Info(ctx context.Context, msg string, fields ...zapcore.Field)
	Warn(ctx context.Context, msg string, fields ...zapcore.Field)
	Error(ctx context.Context, msg string, fields ...zapcore.Field)
	Fatal(ctx context.Context, msg string, fields ...zapcore.Field)
	Panic(ctx context.Context, msg string, fields ...zapcore.Field)
}

type LoggerCtx struct {
	log *zap.Logger
}

func New(log *zap.Logger) *LoggerCtx {
	return &LoggerCtx{log: log}
}

func (l *LoggerCtx) Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.DebugLevel, msg, fields...)
}

func (l *LoggerCtx) Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.InfoLevel, msg, fields...)
}

func (l *LoggerCtx) Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.WarnLevel, msg, fields...)
}

func (l *LoggerCtx) Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.ErrorLevel, msg, fields...)
}

func (l *LoggerCtx) Fatal(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.FatalLevel, msg, fields...)
}

func (l *LoggerCtx) Panic(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.PanicLevel, msg, fields...)
}

func (l *LoggerCtx) logCtx(ctx context.Context, lvl zapcore.Level, m string, f ...zapcore.Field) {
	f = Fields(ctx, f)
	l.log.Log(lvl, msg(ctx, m), f...)
}

