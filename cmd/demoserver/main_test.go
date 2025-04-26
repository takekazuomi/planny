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

package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	todov1 "github.com/takekazu/planny/internal/gen/planny/todo/v1"
	"github.com/takekazu/planny/internal/gen/planny/todo/v1/todov1connect"
)

func TestTodoServer(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.Handle(todov1connect.NewTodoServiceHandler(
		NewTodoServer(),
	))
	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	connectClient := todov1connect.NewTodoServiceClient(
		server.Client(),
		server.URL,
	)
	grpcClient := todov1connect.NewTodoServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)
	clients := []todov1connect.TodoServiceClient{connectClient, grpcClient}

	t.Run("create", func(t *testing.T) {
		for _, client := range clients {
			result, err := client.CreateTodo(context.Background(), connect.NewRequest(&todov1.CreateTodoRequest{
				Title: "Hello",
			}))
			require.NoError(t, err)
			assert.NotEmpty(t, result.Msg.GetTodo().GetName())
		}
	})
}
