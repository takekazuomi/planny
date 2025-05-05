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

// Package: protoconv は、Protocol Buffersのメッセージを変換する
package protoconv

import (
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
)

// ExtractDetails は、与えられたエラーが connect.Error であり、
// その詳細に指定された型 T (proto.Message) の情報が含まれている場合、
// それらをすべて抽出し、スライスとして返す。
// 一致する詳細がない場合や、エラーが connect.Error でない場合は空のスライスを返す。
// https://connectrpc.com/docs/go/errors/#error-details
func ExtractDetails[T proto.Message](err error) []T {
	var connectErr *connect.Error
	var results []T // 結果を格納するスライス。

	if !errors.As(err, &connectErr) {
		return results // connect.Error ではない。
	}

	details := connectErr.Details()
	results = make([]T, 0, len(details))

	for _, detail := range details {
		msg, valueErr := detail.Value()
		if valueErr != nil {
			// 詳細のアンパックに失敗した場合、次の詳細へ。
			continue
		}

		if typedMsg, ok := msg.(T); ok {
			results = append(results, typedMsg) // 型が一致したらスライスに追加。
		}
	}

	return results
}
