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

// Package main is the demo server for the Eliza chatbot. It implements the
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

// CreateTodo implements todov1connect.TodoServiceHandler.
func (t todoServer) CreateTodo(ctx context.Context, req *connect.Request[todov1.CreateTodoRequest]) (*connect.Response[todov1.CreateTodoResponse], error) {
	id := rand.Intn(100)

	todo := &todov1.Todo{
		Name:        "todos/" + strconv.Itoa(id),
		Title:       req.Msg.Title,
		Description: req.Msg.Description,
		Status:      todov1.Todo_STATUS_ACTIVE,
		CreatedAt:   timestamppb.New(time.Now()),
		UpdatedAt:   timestamppb.New(time.Now()),
		DeletedAt:   nil,
	}
	// とりあえず、メモリに保存
	todoData[todo.Name] = todo

	return connect.NewResponse(&todov1.CreateTodoResponse{Todo: todo}), nil
}

// DeleteTodo implements todov1connect.TodoServiceHandler.
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

// GetTodo implements todov1connect.TodoServiceHandler.
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

// ListTodos implements todov1connect.TodoServiceHandler.
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

// UndeleteTodo implements todov1connect.TodoServiceHandler.
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

// UpdateTodo implements todov1connect.TodoServiceHandler.
func (t todoServer) UpdateTodo(context.Context, *connect.Request[todov1.UpdateTodoRequest]) (*connect.Response[todov1.UpdateTodoResponse], error) {
	panic("unimplemented")
}

func NewTodoServer() todov1connect.TodoServiceHandler {
	return &todoServer{}
}
