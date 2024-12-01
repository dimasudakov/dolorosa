package utils

import "golang.org/x/net/context"

// GetTraceID достает trace_id из контекста
func GetTraceID(ctx context.Context) string {
	var traceID string

	// TODO get trace_id from context

	return traceID
}
