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

// Package update は Todo アイテムを更新する機能を提供します
package update

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

// Handler は既存のTodoアイテムを更新します。
func Handler(
	ctx context.Context,
	store store.TodoStore,
	req *connect.Request[todov1.UpdateTodoRequest],
) (*connect.Response[todov1.UpdateTodoResponse], error) {
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

	// 削除されたアイテムは更新できない
	if todo.GetDeletedAt() != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition,
			fmt.Errorf("cannot update deleted todo %q", req.Msg.GetName()))
	}

	// フィールドを更新
	// TODO: パーシャルアップデート
	todo.Title = req.Msg.GetTitle()
	todo.Description = req.Msg.GetDescription()
	todo.UpdatedAt = timestamppb.New(time.Now())

	// レスポンスの作成
	return connect.NewResponse(&todov1.UpdateTodoResponse{
		Todo: todo,
	}), nil
}
