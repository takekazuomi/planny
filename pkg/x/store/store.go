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
	"context"
)

// Store は汎用的なKey-Valueデータストアのインターフェースを定義。
type Store[K comparable, V any] interface {
	// Get はキーに対応する値を取得します。
	Get(ctx context.Context, key K) (V, error)

	// Set はキーに対して値を設定します。
	Set(ctx context.Context, key K, value V) error

	// Delete はキーと対応する値を削除します。
	Delete(ctx context.Context, key K) error

	// List は全ての値を取得します。
	List(ctx context.Context) ([]V, error)

	// Has はキーが存在するか確認します。
	Has(ctx context.Context, key K) bool
}
