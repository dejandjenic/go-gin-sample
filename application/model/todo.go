package model

import (
	"github.com/dejandjenic/go-gin-sample/application/entities"
)

type TodoCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type IdResponse struct {
	Id string
}

type TodoItem struct {
	Name string
	Id   string
}

func (m *TodoCreateRequest) ToEntity() entities.Todo {
	return entities.Todo{
		Name: m.Name,
	}
}

func ToTodoItemSlice(items []entities.Todo) []TodoItem {
	return ToItemSlice[TodoItem, entities.Todo](items, ToTodoItem)
}

func ToItemSlice[K comparable, V any, R any](items []V, converter func(V) R) []R {
	r := []R{}
	for _, i := range items {
		r = append(r, converter(i))
	}
	return r
}

func ToTodoItem(item entities.Todo) TodoItem {
	return TodoItem{
		item.ID,
		item.Name,
	}
}
