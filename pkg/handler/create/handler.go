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

// Package create は Todo アイテムを作成する機能を提供します
package create

import (
	"context"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/oklog/ulid/v2"
	todov1 "github.com/takekazuomi/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazuomi/planny/pkg/store"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Handler は新しいTodoアイテムを作成します。
func Handler(
	ctx context.Context,
	store store.TodoStore,
	req *connect.Request[todov1.CreateTodoRequest],
) (*connect.Response[todov1.CreateTodoResponse], error) {
	// タイムスタンプの作成
	now := time.Now()
	nowProto := timestamppb.New(now)

	// 30日後
	dueTime := now.Add(30 * 24 * time.Hour)
	dueDate := &date.Date{
		Year:  int32(dueTime.Year()),
		Month: int32(dueTime.Month()),
		Day:   int32(dueTime.Day()),
	}

	// Todo リソース名の生成
	name := "todos/" + strings.ToLower(ulid.Make().String())

	// 新しいTodoアイテムの作成
	todo := &todov1.Todo{
		Name:        name,
		Title:       req.Msg.GetTitle(),
		Description: req.Msg.GetDescription(),
		Status:      todov1.TodoStatus_TODO_STATUS_ACTIVE,
		Priority:    0,
		DueDate:     dueDate,
		CreatedAt:   nowProto,
		UpdatedAt:   nowProto,
	}

	// データストアに追加
	err := store.Set(ctx, name, todo)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("failed to create todo: %w", err))
	}

	// レスポンスの作成
	return connect.NewResponse(&todov1.CreateTodoResponse{
		Todo: todo,
	}), nil
}
