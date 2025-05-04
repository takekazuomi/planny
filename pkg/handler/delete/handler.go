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

package delete

import (
	"context"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	todov1 "github.com/takekazu/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazu/planny/pkg/store"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Handler は指定されたリソース名のTodoアイテムを削除します。
func Handler(
	ctx context.Context,
	store store.TodoStore,
	req *connect.Request[todov1.DeleteTodoRequest],
) (*connect.Response[todov1.DeleteTodoResponse], error) {
	// 入力検証
	if req.Msg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("request message is nil"))
	}

	// リソース名からTodoアイテムを取得
	todo, err := store.Get(ctx, req.Msg.GetName())
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound,
			fmt.Errorf("todo with name %q not found", req.Msg.GetName()))
	}

	// 既に削除済みかチェック
	// TODO: AIPを確認
	if todo.DeletedAt != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition,
			fmt.Errorf("todo %q is already deleted", req.Msg.GetName()))
	}

	// 論理削除: 削除日時を設定
	todo.DeletedAt = timestamppb.New(time.Now())
	todo.Status = todov1.TodoStatus_TODO_STATUS_REVOKED

	// レスポンスの作成
	return connect.NewResponse(&todov1.DeleteTodoResponse{
		Todo: todo,
	}), nil
}
