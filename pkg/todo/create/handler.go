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

package create

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	todov1 "github.com/takekazu/planny/pkg/gen/planny/todo/v1"
)

// Handler は新しいTodoアイテムを作成します
func Handler(
	ctx context.Context,
	req *connect.Request[todov1.CreateTodoRequest],
	todoData map[string]*todov1.Todo,
) (*connect.Response[todov1.CreateTodoResponse], error) {
	// 入力検証
	if req.Msg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("request message is nil"))
	}

	// タイムスタンプの作成
	now := time.Now()
	nowProto := timestamppb.New(now)
	// 30日後
	dueDate := timestamppb.New(now.Add(30 * 24 * time.Hour))

	// 新しいTodoアイテムの作成
	todo := &todov1.Todo{
		Name:        fmt.Sprintf("todos/%d", len(todoData)+1),
		Title:       req.Msg.Title,
		Description: req.Msg.Description,
		Status:      todov1.TodoStatus_TODO_STATUS_ACTIVE,
		Priority:    0,
		DueDate:     dueDate,
		CreatedAt:   nowProto,
		UpdatedAt:   nowProto,
	}

	// グローバルなデータストアに追加
	todoData[todo.Name] = todo

	// レスポンスの作成
	return connect.NewResponse(&todov1.CreateTodoResponse{
		Todo: todo,
	}), nil
}
