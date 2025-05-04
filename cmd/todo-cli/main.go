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

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	todov1 "github.com/takekazu/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazu/planny/pkg/gen/planny/todo/v1/todov1connect"
	"google.golang.org/protobuf/proto"
)

func main() {

	conn := &http.Client{
		Timeout: 15 * time.Second,
		// Transport は通常、http.DefaultTransport を基に設定されますが、
		// 必要に応じて TLS 設定などをカスタマイズできます。
		// Transport: http.DefaultTransport,
	}

	client := todov1connect.NewTodoServiceClient(
		conn,
		"localhost:8080",
		connect.WithGRPC(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.CreateTodo(ctx,
		connect.NewRequest(
			&todov1.CreateTodoRequest{
				Title:       "New Todo Item",
				Description: "This is a description of the new todo item.",
			}))

	if err != nil {
		// TodoErrorInfo をすべて抽出。
		todoInfos := ExtractDetails[*todov1.TodoErrorInfo](err)
		if len(todoInfos) > 0 {
			log.Printf("Extracted %d TodoErrorInfo details:", len(todoInfos))
			for i, info := range todoInfos {
				log.Printf("  [%d] Reason=%s, Message=%s", i, info.GetReason(), info.GetMessage())
			}
		}

		// 元のエラーメッセージも出力して終了。
		log.Fatalf("Create failed: %v", err) //nolint:gocritic
	}
	log.Printf("Create response: %v", response)
}

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
