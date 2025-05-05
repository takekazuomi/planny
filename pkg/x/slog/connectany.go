// Copyright 2025 Takekazu Omi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package slog

import (
	"log/slog"
	"net/http"
	"strings"

	"connectrpc.com/connect"
)

// ConnectAnyRequest は connect.AnyRequest から情報を抽出して返却
func ConnectAnyRequest(name string, req connect.AnyRequest) slog.Attr {
	attrs := []any{}

	// メッセージの内容を抽出
	if msg := req.Any(); msg != nil {
		attrs = append(attrs, slog.Any("message", msg))
	}

	// ヘッダー情報を抽出
	if header := extractHeaders(req.Header()); len(header) > 0 {
		attrs = append(attrs, slog.Any("header", header))
	}

	// メソッド情報を抽出
	if method := req.HTTPMethod(); method != "" {
		attrs = append(attrs, slog.String("method", method))
	}

	// 仕様情報を抽出
	spec := req.Spec()
	attrs = append(attrs, ConnectSpec("spec", spec))
	attrs = append(attrs, slog.String("procedure", spec.Procedure))

	// Peer情報を抽出
	peer := req.Peer()
	if peer.Addr != "" {
		attrs = append(attrs, slog.String("peer", peer.Addr))
	}
	if peer.Protocol != "" {
		attrs = append(attrs, slog.String("protocol", peer.Protocol))
	}

	return slog.Group(name, attrs...)
}

// ConnectAnyResponse は connect.AnyResponse から情報を抽出して返却
func ConnectAnyResponse(name string, resp connect.AnyResponse) slog.Attr {
	attrs := []any{}

	// メッセージの内容を抽出
	if msg := resp.Any(); msg != nil {
		attrs = append(attrs, slog.Any("message", msg))
	}

	// ヘッダー情報を抽出
	if header := extractHeaders(resp.Header()); len(header) > 0 {
		attrs = append(attrs, slog.Any("header", header))
	}

	// トレーラー情報を抽出
	if trailer := extractHeaders(resp.Trailer()); len(trailer) > 0 {
		attrs = append(attrs, slog.Any("trailer", trailer))
	}

	// []any型のスライスはslog.Group関数に直接渡せる
	return slog.Group(name, attrs...)
}

// extractHeaders はhttp.Headerからマップを作成して返します
func extractHeaders(header http.Header) map[string][]string {
	if header == nil || len(header) == 0 {
		return nil
	}

	result := make(map[string][]string, len(header))
	for k, v := range header {
		// ヘッダーのキーを正規化 (例: Content-Type -> content-type)
		key := strings.ToLower(k)
		result[key] = v
	}

	return result
}
