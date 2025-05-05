package slog

import (
	"encoding/json"
	"log/slog"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtoMessage(name string, m proto.Message) slog.Attr {
	if m == nil {
		return slog.Any(name, nil)
	}

	// protojson を使って proto.Message を JSON 形式に変換
	marshaler := protojson.MarshalOptions{
		UseProtoNames:   true,  // JSONフィールド名としてprotobufのフィールド名を使用
		EmitUnpopulated: false, // 値が設定されていないフィールドは出力しない
	}

	jsonBytes, err := marshaler.Marshal(m)
	if err != nil {
		// 変換に失敗した場合はそのままのメッセージを返す
		return slog.Any(name, m)
	}

	// JSON バイト列をマップに変換して返す（きれいなログ出力のため）
	var jsonMap map[string]any
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		// JSON からマップへの変換に失敗した場合は文字列として返す
		return slog.String(name, string(jsonBytes))
	}

	return slog.Any(name, jsonMap)
}
