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

package todo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	todov1 "github.com/takekazu/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazu/planny/pkg/gen/planny/todo/v1/todov1connect"
)

func TestCreateTodo(t *testing.T) {
	t.Parallel()

	server, shutdown := StartTodoServer()
	defer shutdown()

	client := todov1connect.NewTodoServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)

	tests := []struct {
		name        string
		req         *todov1.CreateTodoRequest
		wantErr     bool
		checkFields bool
		checkGet    bool
	}{
		{
			name: "基本的なTodo作成",
			req: &todov1.CreateTodoRequest{
				Title:       "重要な会議",
				Description: "経営陣との会議資料の準備",
			},
			wantErr:     false,
			checkFields: true,
			checkGet:    true,
		},
		{
			name: "タイトルのみのTodo作成",
			req: &todov1.CreateTodoRequest{
				Title: "緊急タスク",
			},
			wantErr:     false,
			checkFields: true,
			checkGet:    false,
		},
		{
			name: "説明のみのTodo作成",
			req: &todov1.CreateTodoRequest{
				Description: "後でタイトルを決める",
			},
			wantErr:     false,
			checkFields: true,
			checkGet:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.CreateTodo(context.Background(), connect.NewRequest(tt.req))

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Msg)
			require.NotNil(t, resp.Msg.Todo)

			todo := resp.Msg.Todo

			// 各フィールドが適切に設定されていることを検証
			if tt.checkFields {
				assert.NotEmpty(t, todo.Name, "Todoの名前が設定されていること")
				assert.Equal(t, tt.req.Title, todo.Title, "タイトルが設定されていること")
				assert.Equal(t, tt.req.Description, todo.Description, "説明が設定されていること")
				assert.Equal(t, todov1.TodoStatus_TODO_STATUS_ACTIVE, todo.Status, "ステータスがアクティブに設定されていること")
				assert.NotNil(t, todo.CreatedAt, "作成日時が設定されていること")
				assert.NotNil(t, todo.UpdatedAt, "更新日時が設定されていること")
				assert.Nil(t, todo.DeletedAt, "削除日時が設定されていないこと")
				assert.Equal(t, int32(0), todo.Priority, "優先度がデフォルト値に設定されていること")
				assert.NotNil(t, todo.DueDate, "期限日が初期化されていること")
			}

			// 作成したTodoをGetTodoで取得できることを検証
			if tt.checkGet {
				getResp, err := client.GetTodo(context.Background(), connect.NewRequest(&todov1.GetTodoRequest{
					Name: todo.Name,
				}))

				require.NoError(t, err)
				require.NotNil(t, getResp)
				require.NotNil(t, getResp.Msg)
				require.NotNil(t, getResp.Msg.Todo)

				getTodo := getResp.Msg.Todo
				assert.Equal(t, todo.Name, getTodo.Name, "取得したTodoの名前が一致すること")
				assert.Equal(t, todo.Title, getTodo.Title, "取得したTodoのタイトルが一致すること")
				assert.Equal(t, todo.Description, getTodo.Description, "取得したTodoの説明が一致すること")
			}
		})
	}

	t.Run("複数のTodo作成テスト", func(t *testing.T) {
		// テストケース定義
		todoTests := []struct {
			title       string
			description string
		}{
			{"タスク1", "説明タスク1"},
			{"タスク2", "説明タスク2"},
			{"タスク3", "説明タスク3"},
		}

		// 複数のTodoを作成
		todos := make([]*todov1.Todo, 0, len(todoTests))

		for _, tt := range todoTests {
			req := &todov1.CreateTodoRequest{
				Title:       tt.title,
				Description: tt.description,
			}

			resp, err := client.CreateTodo(context.Background(), connect.NewRequest(req))
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Msg)
			require.NotNil(t, resp.Msg.Todo)

			todos = append(todos, resp.Msg.Todo)

			// 各Todoに対して個別の検証
			assert.Equal(t, tt.title, resp.Msg.Todo.Title)
			assert.Equal(t, tt.description, resp.Msg.Todo.Description)
		}

		// ListTodosを呼び出して、作成したTodoが含まれていることを確認
		listResp, err := client.ListTodos(context.Background(), connect.NewRequest(&todov1.ListTodosRequest{
			ShowDeleted: false,
		}))

		require.NoError(t, err)
		require.NotNil(t, listResp)
		require.NotNil(t, listResp.Msg)

		// 各作成したTodoがリストに含まれていることを確認
		for _, createdTodo := range todos {
			found := false
			for _, listedTodo := range listResp.Msg.Todos {
				if listedTodo.Name == createdTodo.Name {
					found = true
					assert.Equal(t, createdTodo.Title, listedTodo.Title)
					assert.Equal(t, createdTodo.Description, listedTodo.Description)
					break
				}
			}
			assert.True(t, found, "作成したTodo %s がリストに含まれていること", createdTodo.Name)
		}
	})
}

// StartTodoServer は、 Todo service のhttptest.Serverを起動し
// sever と shutdown 関数を返す
func StartTodoServer() (*httptest.Server, func()) {
	mux := http.NewServeMux()
	mux.Handle(todov1connect.NewTodoServiceHandler(
		NewTodoServer(),
	))
	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()

	return server, func() {
		server.Close()
	}
}
