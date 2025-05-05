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

package client

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/takekazuomi/planny/pkg/gen/planny/todo/v1/todov1connect"
)

const (
	userAgent = "planny-client/v1.0.0"
)

// userAgentTransport は、HTTPリクエストにUser-Agentヘッダーを追加するためのTransport
type userAgentTransport struct {
	base      http.RoundTripper
	userAgent string
}

// RoundTrip は、HTTPリクエストを送信する際にUser-Agentヘッダーを追加
func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("user-agent", t.userAgent)
	return t.base.RoundTrip(req)
}

// httpClient は、共通のHTTPクライアント
var httpClient = &http.Client{
	Transport: &userAgentTransport{
		base:      http.DefaultTransport,
		userAgent: userAgent,
	},
}

// NewTodoClient は TodoServiceClient を作成
func NewTodoClient(baseURL string, opts ...connect.ClientOption) todov1connect.TodoServiceClient {
	// http client は、http.DefaultClient を使用
	return todov1connect.NewTodoServiceClient(
		httpClient,
		baseURL,
		opts...,
	)
}
