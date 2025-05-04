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
package server

import (
	"context"

	"connectrpc.com/connect"
	todov1 "github.com/takekazu/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazu/planny/pkg/gen/planny/todo/v1/todov1connect"
	"github.com/takekazu/planny/pkg/handler/create"
	"github.com/takekazu/planny/pkg/handler/delete"
	"github.com/takekazu/planny/pkg/handler/get"
	"github.com/takekazu/planny/pkg/handler/list"
	"github.com/takekazu/planny/pkg/handler/undelete"
	"github.com/takekazu/planny/pkg/handler/update"
	"github.com/takekazu/planny/pkg/store"
)

type todoServer struct {
	Store store.TodoStore
}

var _ todov1connect.TodoServiceHandler = &todoServer{}

// NewTodoServer は新しいTodoServiceHandlerを作成。
// このハンドラーはTodoアイテムの作成、取得、更新、削除、復元などの操作を提供する。
// ConnectフレームワークでTodoサービスを使用するためのエントリーポイントとして使用される。
func NewTodoServer() todov1connect.TodoServiceHandler {
	return &todoServer{
		Store: store.New(),
	}
}

// ハンドラー関数の型を定義。
type handlerFunc[Req any, Resp any] func(context.Context, store.TodoStore, *connect.Request[Req]) (*connect.Response[Resp], error)

// 汎用的なハンドラー呼び出し関数。
func callHandler[Req any, Resp any](
	ctx context.Context, store store.TodoStore, req *connect.Request[Req], handler handlerFunc[Req, Resp],
) (*connect.Response[Resp], error) {
	return handler(ctx, store, req)
}

// CreateTodo は新しいTodoアイテムを作成。
func (s *todoServer) CreateTodo(
	ctx context.Context, req *connect.Request[todov1.CreateTodoRequest],
) (*connect.Response[todov1.CreateTodoResponse], error) {
	return callHandler(ctx, s.Store, req, create.Handler)
}

// GetTodo は指定されたリソース名のTodoアイテムを取得。
func (s *todoServer) GetTodo(
	ctx context.Context, req *connect.Request[todov1.GetTodoRequest],
) (*connect.Response[todov1.GetTodoResponse], error) {
	return callHandler(ctx, s.Store, req, get.Handler)
}

// UpdateTodo は既存のTodoアイテムを更新。
func (s *todoServer) UpdateTodo(
	ctx context.Context, req *connect.Request[todov1.UpdateTodoRequest],
) (*connect.Response[todov1.UpdateTodoResponse], error) {
	return callHandler(ctx, s.Store, req, update.Handler)
}

// ListTodos はTodoアイテムのリストを取得。
func (s *todoServer) ListTodos(
	ctx context.Context, req *connect.Request[todov1.ListTodosRequest],
) (*connect.Response[todov1.ListTodosResponse], error) {
	return callHandler(ctx, s.Store, req, list.Handler)
}

// DeleteTodo は指定されたリソース名のTodoアイテムを削除。
func (s *todoServer) DeleteTodo(
	ctx context.Context, req *connect.Request[todov1.DeleteTodoRequest],
) (*connect.Response[todov1.DeleteTodoResponse], error) {
	return callHandler(ctx, s.Store, req, delete.Handler)
}

// UndeleteTodo は削除されたTodoアイテムを復元。
func (s *todoServer) UndeleteTodo(ctx context.Context, req *connect.Request[todov1.UndeleteTodoRequest],
) (*connect.Response[todov1.UndeleteTodoResponse], error) {
	return callHandler(ctx, s.Store, req, undelete.Handler)
}
