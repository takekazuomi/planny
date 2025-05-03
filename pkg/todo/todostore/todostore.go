package todostore

import (
	todov1 "github.com/takekazu/planny/pkg/gen/planny/todo/v1"
	"github.com/takekazu/planny/pkg/store"
)

type TodoStore = store.KeyValueStore[string, *todov1.Todo]
