package logging

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"pay/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ctxKey string

const (
	keyLogger  ctxKey = "logger"
	keyTraceID ctxKey = "trace_id"
)

func FromGin(c *gin.Context) *zap.Logger {
	if c == nil {
		return L()
	}
	if v, ok := c.Get(string(keyLogger)); ok {
		if l, ok := v.(*zap.Logger); ok && l != nil {
			return l
		}
	}
	return L()
}

func TraceIDFromGin(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if v, ok := c.Get(string(keyTraceID)); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(keyTraceID).(string); ok {
		return v
	}
	return ""
}

func Middleware(base *zap.Logger, cfg config.TraceConfig) gin.HandlerFunc {
	if base == nil {
		base = L()
	}
	header := cfg.Header
	if header == "" {
		header = "X-Trace-Id"
	}

	return func(c *gin.Context) {
		traceID := ""
		if cfg.Enabled {
			traceID = extractTraceID(c, header)
			if traceID == "" {
				traceID = newTraceID()
			}
			c.Set(string(keyTraceID), traceID)
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), keyTraceID, traceID))
			c.Writer.Header().Set(header, traceID)
		}

		reqLogger := base
		if traceID != "" {
			reqLogger = reqLogger.With(zap.String("trace_id", traceID))
		}
		c.Set(string(keyLogger), reqLogger)

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		}
		if c.Errors != nil && len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		status := c.Writer.Status()
		switch {
		case status >= 500:
			reqLogger.Error("http_request", fields...)
		case status >= 400:
			reqLogger.Warn("http_request", fields...)
		default:
			reqLogger.Info("http_request", fields...)
		}
	}
}

func extractTraceID(c *gin.Context, header string) string {
	if c == nil {
		return ""
	}
	if v := c.GetHeader(header); v != "" {
		return sanitizeTraceID(v)
	}
	if v := c.GetHeader("traceparent"); v != "" {
		parts := strings.Split(v, "-")
		if len(parts) >= 3 {
			return sanitizeTraceID(parts[1])
		}
	}
	return ""
}

func sanitizeTraceID(v string) string {
	s := strings.TrimSpace(v)
	if len(s) > 64 {
		s = s[:64]
	}
	return s
}

func newTraceID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
