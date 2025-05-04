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

// Package list は Todo アイテムを取得する機能を提供します
package list

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	todov1 "github.com/takekazu/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazu/planny/pkg/store"
)

// Handler はTodoアイテムのリストを取得します。
func Handler(
	ctx context.Context,
	store store.TodoStore,
	req *connect.Request[todov1.ListTodosRequest],
) (*connect.Response[todov1.ListTodosResponse], error) {
	// 入力検証
	if req.Msg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("request message is nil"))
	}

	// todoDataから条件に合うアイテムを抽出
	todos, err := store.List(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to list todos: %w", err))
	}

	// 結果を格納するスライス
	results := make([]*todov1.Todo, 0, len(todos))

	// フィルタリング
	for _, todo := range todos {
		// 削除済みのアイテムの表示/非表示を制御
		if todo.GetDeletedAt() != nil && !req.Msg.GetShowDeleted() {
			continue
		}
		results = append(results, todo)
	}

	// レスポンスの作成
	return connect.NewResponse(&todov1.ListTodosResponse{
		Todos: results,
	}), nil
}
