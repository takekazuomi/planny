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
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// MethodDescriptorToMap は protoreflect.MethodDescriptor の重要な情報を
// サイクル参照を回避して安全にマップに変換。
func MethodDescriptorToMap(md protoreflect.MethodDescriptor) map[string]interface{} { //nolint:cyclop,varnamelen
	if md == nil {
		return nil
	}

	result := map[string]any{
		"name":             string(md.Name()),
		"full_name":        string(md.FullName()),
		"is_streaming":     md.IsStreamingClient() || md.IsStreamingServer(),
		"is_client_stream": md.IsStreamingClient(),
		"is_server_stream": md.IsStreamingServer(),
	}

	// 入力メッセージの型情報
	if md.Input() != nil {
		result["input_type"] = string(md.Input().FullName())
	}

	// 出力メッセージの型情報
	if md.Output() != nil {
		result["output_type"] = string(md.Output().FullName())
	}

	// コメントがあれば追加（安全なアクセス）
	if md.ParentFile() != nil {
		parentIdx := md.Index()
		locations := md.ParentFile().SourceLocations()

		// parentIdxの有効値確認
		if parentIdx >= 0 && parentIdx < locations.Len() {
			location := locations.Get(parentIdx)
			if location.Path != nil && location.LeadingComments != "" {
				result["comments"] = location.LeadingComments
			}
		}
	}

	// 親サービス情報の追加
	if parent := md.Parent(); parent != nil {
		if service, ok := parent.(protoreflect.ServiceDescriptor); ok {
			result["service"] = string(service.FullName())
		}
	}

	return result
}

// MethodDescriptor は protoreflect.MethodDescriptor を slog.Attr として安全に返却。
func MethodDescriptor(name string, md protoreflect.MethodDescriptor) slog.Attr {
	return slog.Any(name, MethodDescriptorToMap(md))
}

// ConnectSpecSchema は connect.Spec.Schema から安全に MethodDescriptor を抽出。
func ConnectSpecSchema(name string, schema any) slog.Attr {
	if schema == nil {
		return slog.Any(name, nil)
	}

	// Schema が MethodDescriptor であればそれを使用
	if md, ok := schema.(protoreflect.MethodDescriptor); ok {
		return MethodDescriptor(name, md)
	}

	// それ以外はtype情報のみを返却
	return slog.String(name+"_type", fmt.Sprintf("%T", schema))
}

// ConnectSpec は connect.Spec から情報を抽出して返却。
func ConnectSpec(name string, spec connect.Spec) slog.Attr {
	// マップに変換して安全に扱う
	result := map[string]interface{}{
		"procedure":         spec.Procedure,
		"stream_type":       spec.StreamType.String(),
		"is_client":         spec.IsClient,
		"idempotency_level": spec.IdempotencyLevel.String(),
	}

	// Schemaフィールドの処理（通常はMethodDescriptor）
	schemaAttr := ConnectSpecSchema("schema", spec.Schema)
	if schemaAttr.Value.Kind() != slog.KindAny || schemaAttr.Value.Any() != nil {
		// スキーマが存在する場合のみ含める
		result["schema"] = schemaAttr.Value.Any()
	}

	return slog.Any(name, result)
}
