package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logCtx struct {
	MessageID   string
	HandlerName string
	MsgKey      string
	Message     string
}

type keyType int

const key keyType = iota

func Fields(ctx context.Context, fields []zapcore.Field) []zapcore.Field {
	l, ok := ctx.Value(key).(logCtx)
	if !ok {
		return fields
	}

	if l.MessageID != "" {
		fields = append(fields, zap.String("message_id", l.Message))
	}

	if l.HandlerName != "" {
		fields = append(fields, zap.String("handler_name", l.HandlerName))
	}

	if l.MsgKey != "" {
		fields = append(fields, zap.String("msg_key", l.MsgKey))
	}

	if l.Message != "" {
		fields = append(fields, zap.String("message", l.Message))
	}

	return fields
}

func msg(ctx context.Context, msg string) string {
	if l, ok := ctx.Value(key).(logCtx); ok {
		if l.HandlerName != "" {
			return fmt.Sprintf("%s: %s", l.HandlerName, msg)
		}
	}

	return msg
}

func WithMessageID(ctx context.Context, messageID string) context.Context {
	if l, ok := ctx.Value(key).(logCtx); ok {
		l.MessageID = messageID

		return context.WithValue(ctx, key, l)
	}

	return context.WithValue(ctx, key, logCtx{MessageID: messageID})
}

func WithHandlerName(ctx context.Context, handlerName string) context.Context {
	if l, ok := ctx.Value(key).(logCtx); ok {
		l.HandlerName = handlerName

		return context.WithValue(ctx, key, l)
	}

	return context.WithValue(ctx, key, logCtx{HandlerName: handlerName})
}

func WithMsgKey(ctx context.Context, msgKey string) context.Context {
	if l, ok := ctx.Value(key).(logCtx); ok {
		l.MsgKey = msgKey

		return context.WithValue(ctx, key, l)
	}

	return context.WithValue(ctx, key, logCtx{MsgKey: msgKey})
}

func WithMessage(ctx context.Context, msg string) context.Context {
	if l, ok := ctx.Value(key).(logCtx); ok {
		l.MsgKey = msg

		return context.WithValue(ctx, key, l)
	}

	return context.WithValue(ctx, key, logCtx{Message: msg})
}

