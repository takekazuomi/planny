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

package undelete

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	todov1 "github.com/takekazu/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazu/planny/pkg/todo/store"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Handler は削除されたTodoアイテムを復元します
func Handler(
	ctx context.Context,
	store store.TodoStore,
	req *connect.Request[todov1.UndeleteTodoRequest],
) (*connect.Response[todov1.UndeleteTodoResponse], error) {
	// 入力検証
	if req.Msg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("request message is nil"))
	}

	// リソース名からTodoアイテムを取得
	todo, err := store.Get(ctx, req.Msg.Name)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound,
			fmt.Errorf("todo with name %q not found", req.Msg.Name))
	}

	// 削除されていない場合はエラー
	if todo.DeletedAt == nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition,
			fmt.Errorf("todo %q is not deleted", req.Msg.Name))
	}

	// 削除フラグをクリア
	todo.DeletedAt = nil
	todo.Status = todov1.TodoStatus_TODO_STATUS_ACTIVE
	todo.UpdatedAt = timestamppb.New(time.Now())

	// レスポンスの作成
	return connect.NewResponse(&todov1.UndeleteTodoResponse{
		Todo: todo,
	}), nil
}
