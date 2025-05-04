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

package store

import (
	todov1 "github.com/takekazu/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazu/planny/pkg/x/store"
	"github.com/takekazu/planny/pkg/x/store/mapstore"
)

// TodoStore はTodoアイテムを管理するためのインターフェース。
type TodoStore interface {
	store.Store[string, *todov1.Todo]
}

// store はTodoStoreの実装。
type todoStore struct {
	// 構造体埋め込みによりKeyValueStoreの実装を継承
	store.Store[string, *todov1.Todo]
}

var _ store.Store[string, *todov1.Todo] = (*todoStore)(nil)

// New は新しいTodoStoreのインスタンスを作成。
func New() TodoStore {
	return &todoStore{
		Store: mapstore.New[string, *todov1.Todo](),
	}
}
