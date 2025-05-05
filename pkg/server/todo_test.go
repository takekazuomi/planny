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

package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	todov1 "github.com/takekazuomi/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazuomi/planny/pkg/gen/planny/todo/v1/todov1connect"
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resp, err := client.CreateTodo(context.Background(), connect.NewRequest(testCase.req))

			if testCase.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Msg)
			require.NotNil(t, resp.Msg.GetTodo())

			todo := resp.Msg.GetTodo()

			// 各フィールドが適切に設定されていることを検証
			if testCase.checkFields {
				assert.NotEmpty(t, todo.GetName(), "Todoの名前が設定されていること")
				assert.Equal(t, testCase.req.GetTitle(), todo.GetTitle(), "タイトルが設定されていること")
				assert.Equal(t, testCase.req.GetDescription(), todo.GetDescription(), "説明が設定されていること")
				assert.Equal(t, todov1.TodoStatus_TODO_STATUS_ACTIVE, todo.GetStatus(), "ステータスがアクティブに設定されていること")
				assert.NotNil(t, todo.GetCreatedAt(), "作成日時が設定されていること")
				assert.NotNil(t, todo.GetUpdatedAt(), "更新日時が設定されていること")
				assert.Nil(t, todo.GetDeletedAt(), "削除日時が設定されていないこと")
				assert.Equal(t, int32(0), todo.GetPriority(), "優先度がデフォルト値に設定されていること")
				assert.NotNil(t, todo.GetDueDate(), "期限日が初期化されていること")
			}

			// 作成したTodoをGetTodoで取得できることを検証
			if testCase.checkGet {
				getResp, err := client.GetTodo(context.Background(), connect.NewRequest(&todov1.GetTodoRequest{
					Name: todo.GetName(),
				}))

				require.NoError(t, err)
				require.NotNil(t, getResp)
				require.NotNil(t, getResp.Msg)
				require.NotNil(t, getResp.Msg.GetTodo())

				getTodo := getResp.Msg.GetTodo()
				assert.Equal(t, todo.GetName(), getTodo.GetName(), "取得したTodoの名前が一致すること")
				assert.Equal(t, todo.GetTitle(), getTodo.GetTitle(), "取得したTodoのタイトルが一致すること")
				assert.Equal(t, todo.GetDescription(), getTodo.GetDescription(), "取得したTodoの説明が一致すること")
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

		for _, todoTest := range todoTests {
			req := &todov1.CreateTodoRequest{
				Title:       todoTest.title,
				Description: todoTest.description,
			}

			resp, err := client.CreateTodo(context.Background(), connect.NewRequest(req))
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Msg)
			require.NotNil(t, resp.Msg.GetTodo())

			todos = append(todos, resp.Msg.GetTodo())

			// 各Todoに対して個別の検証
			assert.Equal(t, todoTest.title, resp.Msg.GetTodo().GetTitle())
			assert.Equal(t, todoTest.description, resp.Msg.GetTodo().GetDescription())
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
			for _, listedTodo := range listResp.Msg.GetTodos() {
				if listedTodo.GetName() == createdTodo.GetName() {
					found = true
					assert.Equal(t, createdTodo.GetTitle(), listedTodo.GetTitle())
					assert.Equal(t, createdTodo.GetDescription(), listedTodo.GetDescription())
					break
				}
			}
			assert.True(t, found, "作成したTodo %s がリストに含まれていること", createdTodo.GetName())
		}
	})
}

// StartTodoServer は、 Todo service のhttptest.Serverを起動し server と shutdown 関数を返す。
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
