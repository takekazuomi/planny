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

// Package: mapstore は、KeyValueStore のin-memory generic map実装
package mapstore

import (
	"context"
	"fmt"
	"sync"

	"github.com/takekazu/planny/pkg/x/store"
)

// MapStore はKeyValueStoreインターフェースを実装する、メモリ内のmapベースストア
type MapStore[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

// MapStoreはKeyValueStoreインターフェースを実装
var _ store.Store[string, any] = (*MapStore[string, any])(nil)

// initialCapacity は固定
const initialCapacity = 1000

// New は新しいMapStoreインスタンスの作成
func New[K comparable, V any]() *MapStore[K, V] {
	return &MapStore[K, V]{
		data: make(map[K]V, initialCapacity),
	}
}

// Get はキーに対応する値の取得
func (s *MapStore[K, V]) Get(ctx context.Context, key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	if !exists {
		var zero V
		return zero, fmt.Errorf("key not found: %v", key)
	}

	return value, nil
}

// Set はキーに対する値の設定
func (s *MapStore[K, V]) Set(ctx context.Context, key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return nil
}

// Delete はキーと対応する値の削除
func (s *MapStore[K, V]) Delete(ctx context.Context, key K) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	return nil
}

// List は全ての値の取得
func (s *MapStore[K, V]) List(ctx context.Context) ([]V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	values := make([]V, 0, len(s.data))
	for _, v := range s.data {
		values = append(values, v)
	}

	return values, nil
}

// Has はキーの存在確認
func (s *MapStore[K, V]) Has(ctx context.Context, key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.data[key]
	return exists
}
