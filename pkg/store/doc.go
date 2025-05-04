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

// Package store はTodoアプリケーションのためのデータ永続化レイヤーの提供。
//
// このパッケージでは、Todoアイテムを管理するための抽象化ストレージインターフェース
// と、実装の定義。パッケージの主な責務は
// 以下の通り:
//
//   - TodoStoreインターフェース: Todoアイテムの CRUD 操作の定義。
//   - 複数のストレージバックエンドの切り替えのIF。
//   - トランザクション管理: アトミックな操作のサポート。
//   - エラーハンドリング: 一貫性のある方法でのストレージ操作エラーの処理。
//   - 現状は、メモリ内ストレージの実装を提供。
package store
