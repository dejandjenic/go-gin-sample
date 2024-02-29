package tests

import (
	"net/http"
	"testing"

	"github.com/dejandjenic/go-gin-sample/application/model"
	"github.com/dejandjenic/go-gin-sample/application/setup"
)

var auth = "Bearer " + "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJHclhadkRIaFBnU1JHRzBGU0dLT1RER3pualE3cWhfUlhPUlZTZ28wMURVIn0.eyJleHAiOjE3MDkyMzcwNDQsImlhdCI6MTcwOTIwMTA0NSwiYXV0aF90aW1lIjoxNzA5MjAxMDQ0LCJqdGkiOiI1ODgwZTcxMC1hMTNjLTQ2NGMtODk4ZC0yNTZmZmRiZmI2MjEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwODEvcmVhbG1zL0RFTU9SRUFMTSIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiIzNWUyMTk2ZC1hMTcxLTQyNjgtOTYxZi1mMGNhZDM3YWI3NzEiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhY2NvdW50LWNvbnNvbGUiLCJub25jZSI6IjliNzE2M2YyLWQ4NDgtNDVkNy04NzNjLTY4ODE1YTNkMTRjYyIsInNlc3Npb25fc3RhdGUiOiIyMmY0NzBmOS1iYWFkLTQwMTQtOTk2Ni0wNmU0MWVlMmYyY2IiLCJhY3IiOiIxIiwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyJdfX0sInNjb3BlIjoib3BlbmlkIGVtYWlsIHByb2ZpbGUiLCJzaWQiOiIyMmY0NzBmOS1iYWFkLTQwMTQtOTk2Ni0wNmU0MWVlMmYyY2IiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiZGVqYW4ifQ.K6mrGoo1Y7cxdxA76UTE896bAzvfu--9D1GX9uJIs5MKs4gk90vfA1qJY_qRuQKW20evxWA_KDwvmISfAVhQSnL5r7ACkGgnGSzKXFU_b4nj5_SvezKD3fyqu4TLgidyNyA7A-MiBisqf09ZUhSc9AkfHJiJpbGhBTF3SLkiDKjviX-axyvZozdoXHn76LO6-33lNrdQzLdui8Q-tqBctNz4WG88u4nI_3lechyFWYsThJp2Jy_z1EZHLX5GwbpYl6MT-Jvkq_wzM3jYWPDTPJghBxDlHEj4_z5RFzuMqE_C4_QuI-vPHgWjiEzg0SnRNBtFUFjfKZu7z_F2wZziCg"

var _, url = setupMockServer()
var handler, _ = setup.GinHandler(url)

func TestUnAuthorized(t *testing.T) {
	e := getExpect(t, handler)

	e.GET("/api/v1/ping").
		Expect().
		Status(http.StatusUnauthorized)
}

func TestPong(t *testing.T) {
	e := getExpect(t, handler)
	e.
		GET("/api/v1/ping").
		WithHeader("Authorization", auth).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		HasValue("message", "pong")
}

func TestFlow(t *testing.T) {
	e := getExpect(t, handler)
	e.GET("/api/v1/todos").
		WithHeader("Authorization", auth).
		Expect().
		Status(http.StatusOK).
		JSON().
		Array().
		Length().
		IsEqual(0)

	var response model.IdResponse

	e.POST("/api/v1/todos").
		WithHeader("Authorization", auth).
		WithJSON(model.TodoCreateRequest{
			Name: "test",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&response)

	e.GET("/api/v1/todos").
		WithHeader("Authorization", auth).
		Expect().
		Status(http.StatusOK).
		JSON().
		Array().
		Length().
		IsEqual(1)

	e.GET("/api/v1/todos/"+response.Id).
		WithHeader("Authorization", auth).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		HasValue("Name", "test")

	e.PUT("/api/v1/todos/"+response.Id).
		WithHeader("Authorization", auth).
		WithJSON(model.TodoCreateRequest{
			Name: "xxx",
		}).
		Expect().
		Status(http.StatusOK)

	e.GET("/api/v1/todos/"+response.Id).
		WithHeader("Authorization", auth).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		HasValue("Name", "xxx")

	e.DELETE("/api/v1/todos/"+response.Id).
		WithHeader("Authorization", auth).
		Expect().
		Status(http.StatusOK)

	e.GET("/api/v1/todos/"+response.Id).
		WithHeader("Authorization", auth).
		Expect().
		Status(http.StatusNotFound)

	e.GET("/api/v1/todos").
		WithHeader("Authorization", auth).
		Expect().
		Status(http.StatusOK).
		JSON().
		Array().
		Length().
		IsEqual(0)
}
