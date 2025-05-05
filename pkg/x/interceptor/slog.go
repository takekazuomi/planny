// Package interceptor は、Connect RPCのインターセプターを提供する。
package interceptor

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	xslog "github.com/takekazuomi/planny/pkg/x/slog"
)

// NewSlogInterceptor は、リクエストとレスポンスをログに記録するインターセプターを作成
func NewSlogInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			procedure := req.Spec().Procedure

			attrs := []slog.Attr{}
			attrs = append(attrs, xslog.ConnectAnyRequest("req", req))

			logger.LogAttrs(ctx, slog.LevelInfo, "grpc call started", attrs...)

			// gRPCリクエストを処理
			resp, err := next(ctx, req)

			attrs = []slog.Attr{
				slog.String("procedure", procedure),
				slog.Float64("duration_ms", float64(time.Since(start).Nanoseconds())/1000000.0),
			}
			if err != nil {
				attrs = append(attrs, slog.Any("error", err))
			}
			attrs = append(attrs, xslog.ConnectAnyResponse("resp", resp))

			logger.LogAttrs(ctx, slog.LevelInfo, "grpc call finished", attrs...)

			return resp, err
		}
	}
}

// TestJSONMarshal はオブジェクトがJSONマーシャリング可能か検証します
func TestJSONMarshal(obj any, name string) {
	_, err := json.Marshal(obj)
	if err != nil {
		log.Printf("JSON marshaling error for %s: %v", name, err)
	} else {
		log.Printf("%s can be marshaled to JSON successfully", name)
	}
}
