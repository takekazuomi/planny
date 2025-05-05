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

// Package opts は、コマンドラインオプションを定義。
package opts

var (
	// APIEndpoint は、APIエンドポイントのURL。
	APIEndpoint = "http://localhost:51051"
	// Title は、TODOのタイトル。
	Title = ""
	// Description は、TODOの説明。
	Description = ""
	// DueDate は、TODOの予定。
	DueDate = ""
	// Priority は、TODOの優先度。
	Priority int32
)
