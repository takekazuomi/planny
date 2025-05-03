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

package list

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	todov1 "github.com/takekazu/planny/pkg/gen/planny/todo/v1"
)

// Handler はTodoアイテムのリストを取得します
func Handler(
	ctx context.Context,
	req *connect.Request[todov1.ListTodosRequest],
	todoData map[string]*todov1.Todo,
) (*connect.Response[todov1.ListTodosResponse], error) {
	// 入力検証
	if req.Msg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("request message is nil"))
	}

	// 結果を格納するスライス
	todos := make([]*todov1.Todo, 0, len(todoData))

	// todoDataから条件に合うアイテムを抽出
	for _, todo := range todoData {
		// 削除済みのアイテムの表示/非表示を制御
		if todo.DeletedAt != nil && !req.Msg.ShowDeleted {
			continue
		}
		todos = append(todos, todo)
	}

	// レスポンスの作成
	return connect.NewResponse(&todov1.ListTodosResponse{
		Todos: todos,
	}), nil
}
