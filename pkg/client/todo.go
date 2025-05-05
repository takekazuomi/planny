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
