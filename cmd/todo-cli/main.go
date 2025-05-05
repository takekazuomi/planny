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

// Package main はTodoサービスのクライアントコマンドラインツール
package main

import (
	"github.com/goaux/results"
	"github.com/goaux/slog/logger"
	"github.com/takekazuomi/planny/pkg/cmd"
)

var log = results.Must1(logger.NewName("planny-cli")) //nolint:unused,gochecknoglobals

func main() {
	cmd.Execute()
}
