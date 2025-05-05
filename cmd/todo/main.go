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

// Package main implements the gRPC server for the Todo service.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"

	"github.com/goaux/results"
	"github.com/goaux/slog/logger"
	"github.com/takekazuomi/planny/pkg/gen/planny/todo/v1/todov1connect"
	"github.com/takekazuomi/planny/pkg/server"
	"github.com/takekazuomi/planny/pkg/x/interceptor"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var log = results.Must1(logger.NewName("planny"))

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()

	mux.Handle(todov1connect.NewTodoServiceHandler(
		server.NewTodoServer(),
		connect.WithInterceptors(interceptor.NewSlogInterceptor(log)),
	))

	mux.Handle(grpchealth.NewHandler(
		grpchealth.NewStaticChecker(todov1connect.TodoServiceName),
	))

	mux.Handle(grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(todov1connect.TodoServiceName),
	))

	mux.Handle(grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(todov1connect.TodoServiceName),
	))

	addr := "localhost:51051"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	srv := &http.Server{
		Addr: addr,
		Handler: h2c.NewHandler(
			mux,
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.ErrorContext(ctx, "HTTP listen and serve", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-signals

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.ErrorContext(ctx, "HTTP shutdown", slog.Any("error", err))
	}
}
