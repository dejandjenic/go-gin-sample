package handlers

import "github.com/dejandjenic/go-gin-sample/application"

type Handler application.Application

func ToHandler(a *application.Application) Handler {
	return Handler{
		TodoRepository: a.TodoRepository,
		Configuration:  a.Configuration,
	}
}
