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

// Package cmd は、コマンドラインインターフェースを提供するパッケージ。
package cmd

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
	"github.com/takekazuomi/planny/pkg/client"
	"github.com/takekazuomi/planny/pkg/cmd/opts"
	todov1 "github.com/takekazuomi/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazuomi/planny/pkg/x/protoconv"
)

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: run,
}

func init() {
	rootCmd.AddCommand(createCmd)

	flags := createCmd.Flags()

	flags.StringVarP(&opts.Title, "title", "t", opts.Title, "Todo title")
	flags.StringVarP(&opts.Description, "description", "d", opts.Description, "Todo description")

	opts.DueDate = time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02")
	flags.StringVarP(&opts.DueDate, "due-date", "D", opts.DueDate, "Todo due date. Format: YYYY-MM-DD")

	flags.Int32VarP(&opts.Priority, "priority", "p", 0, "Todo priority")
}

func run(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	log = log.With(slog.String("cmd", "create"))

	client := client.NewTodoClient(opts.APIEndpoint, connect.WithGRPC())

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dueTime, err := time.Parse("2006-01-02", opts.DueDate)
	if err != nil {
		// 日付のフォーマットが無効な場合、エラーメッセージを表示して終了。
		return err
	}

	dueDate := protoconv.TimeToDate(dueTime)

	response, err := client.CreateTodo(ctx,
		connect.NewRequest(
			&todov1.CreateTodoRequest{
				Title:       opts.Title,
				Description: opts.Description,
				DueDate:     dueDate,
				Priority:    opts.Priority,
			}))

	if err != nil {
		// TodoErrorInfo をすべて抽出。
		todoInfos := protoconv.ExtractDetails[*todov1.TodoErrorInfo](err)
		if len(todoInfos) > 0 {
			log.InfoContext(ctx, "ExtractDetails", slog.Any("TodoErrorInfo", todoInfos))
		}

		// 元のエラーメッセージも出力
		log.ErrorContext(ctx, "failed", slog.Any("error", err))
		return err
	}

	log.Info("done", slog.Any("response", response))
	return nil
}
