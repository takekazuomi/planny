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

// Package todo はTodoサービスを実装するパッケージです。
// TodoServiceHandlerインターフェースを実装し、Todoアイテムの管理機能を提供します。
package todo

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"connectrpc.com/connect"
	todov1 "github.com/takekazu/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazu/planny/pkg/gen/planny/todo/v1/todov1connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type todoServer struct {
}

var _ todov1connect.TodoServiceHandler = todoServer{}

var todoData map[string]*todov1.Todo = make(map[string]*todov1.Todo, 100)

// CreateTodo は新しいTodoアイテムを作成します。
// リクエストからタイトルと説明を取得し、ランダムなIDを生成して新しいTodoアイテムを作成します。
// 作成されたTodoアイテムはメモリ内に保存されます。
func (t todoServer) CreateTodo(ctx context.Context, req *connect.Request[todov1.CreateTodoRequest]) (*connect.Response[todov1.CreateTodoResponse], error) {
	id := rand.Intn(100)

	todo := &todov1.Todo{
		Name:        "todos/" + strconv.Itoa(id),
		Title:       req.Msg.Title,
		Description: req.Msg.Description,
		Status:      todov1.TodoStatus_TODO_STATUS_ACTIVE,
		CreatedAt:   timestamppb.New(time.Now()),
		UpdatedAt:   timestamppb.New(time.Now()),
		DeletedAt:   nil,
		Priority:    0,
		DueDate:     &timestamppb.Timestamp{},
	}
	// とりあえず、メモリに保存
	todoData[todo.Name] = todo

	return connect.NewResponse(&todov1.CreateTodoResponse{Todo: todo}), nil
}

// DeleteTodo は指定されたTodoアイテムを論理削除します。
// 指定されたリソース名のTodoアイテムに削除フラグ（削除日時）を設定します。
// 実際にはデータはメモリから削除されません。
func (t todoServer) DeleteTodo(ctx context.Context, req *connect.Request[todov1.DeleteTodoRequest]) (*connect.Response[todov1.DeleteTodoResponse], error) {
	todoName := req.Msg.Name
	todo, ok := todoData[todoName]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("todo not found"))
	}
	// メモリ上から消さずに、削除フラグを立てる
	todo.DeletedAt = timestamppb.New(time.Now())
	// 保存
	todoData[todoName] = todo
	return connect.NewResponse(&todov1.DeleteTodoResponse{Todo: todo}), nil
}

// GetTodo は指定されたリソース名のTodoアイテムを取得します。
// 削除フラグが立っているTodoアイテムも取得できます。
func (t todoServer) GetTodo(ctx context.Context, req *connect.Request[todov1.GetTodoRequest]) (*connect.Response[todov1.GetTodoResponse], error) {
	todoID := req.Msg.Name
	todo, ok := todoData[todoID]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("todo not found"))
	}
	// 取得したTodoを返す
	// 削除フラグが立っている場合も、そのまま返す
	return connect.NewResponse(&todov1.GetTodoResponse{Todo: todo}), nil
}

// ListTodos はTodoアイテムのリストを取得します。
// showDeletedパラメータがtrueの場合は、削除されたTodoアイテムも含めて取得します。
func (t todoServer) ListTodos(ctx context.Context, req *connect.Request[todov1.ListTodosRequest]) (*connect.Response[todov1.ListTodosResponse], error) {
	// 全てのTodoを取得する
	todos := make([]*todov1.Todo, 0, len(todoData))
	for _, todo := range todoData {
		// 削除フラグが立っている場合は、取得しない
		if req.Msg.GetShowDeleted() || todo.DeletedAt == nil {
			todos = append(todos, todo)
		}
	}
	// 取得したTodoを返す
	return connect.NewResponse(&todov1.ListTodosResponse{Todos: todos}), nil
}

// UndeleteTodo は削除済みのTodoアイテムを復元します。
// 削除フラグ（DeletedAt）を解除して、Todoアイテムを再度アクティブにします。
// 削除されていないTodoアイテムに対して呼び出した場合はエラーを返します。
func (t todoServer) UndeleteTodo(ctx context.Context, req *connect.Request[todov1.UndeleteTodoRequest]) (*connect.Response[todov1.UndeleteTodoResponse], error) {
	// UndeleteTodoは、削除フラグを立てたTodoを復活させる
	var todo *todov1.Todo
	var ok bool
	if todo, ok = todoData[req.Msg.GetName()]; !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("todo not found"))
	}

	// 削除フラグが立っていない場合は、エラーを返す
	if todo.DeletedAt == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("todo not found"))
	}
	// 削除フラグを消す
	todo.DeletedAt = nil

	return connect.NewResponse(&todov1.UndeleteTodoResponse{Todo: todo}), nil
}

// UpdateTodo は既存のTodoアイテムを更新します。
// 現在は未実装です。
func (t todoServer) UpdateTodo(context.Context, *connect.Request[todov1.UpdateTodoRequest]) (*connect.Response[todov1.UpdateTodoResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

// NewTodoServer は新しいTodoServiceHandlerを作成します。
// このハンドラーはTodoアイテムの作成、取得、更新、削除、復元などの操作を提供します。
// ConnectフレームワークでTodoサービスを使用するためのエントリーポイントとして使用されます。
func NewTodoServer() todov1connect.TodoServiceHandler {
	return &todoServer{}
}
