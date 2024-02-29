package handlers

import (
	"fmt"

	"github.com/dejandjenic/go-gin-sample/application/middleware"
	"github.com/dejandjenic/go-gin-sample/application/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @Summary create todo
// @Schemes
// @Description create todo
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} model.IdResponse
// @Failure 400
// @Router /todos [post]
func (h *Handler) CreateTodo(g *gin.Context) {
	var r model.TodoCreateRequest
	berr := g.ShouldBindJSON(&r)

	if berr != nil {
		g.JSON(400, gin.H{
			"error": fmt.Sprintf("An error %s", berr.Error()),
		})
		return
	}

	id, err := h.TodoRepository.CreateTodo(g.Request.Context(), r.ToEntity())

	if err != nil {
		g.JSON(400, gin.H{
			"error": "an error occurred",
		})
		return
	}

	g.JSON(200,
		model.IdResponse{
			Id: id,
		},
	)
}

// @Summary update todo
// @Schemes
// @Description update todo
// @Tags example
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Router /todos/{id} [put]
// @Param id path string true "todo id"
func (h *Handler) UpdateTodo(g *gin.Context) {
	id := g.Param("id")

	var r model.TodoCreateRequest
	berr := g.ShouldBindJSON(&r)

	if berr != nil {
		g.JSON(400, gin.H{
			"error": fmt.Sprintf("An error %s", berr.Error()),
		})
		return
	}

	h.TodoRepository.UpdateTodo(g.Request.Context(), id, r.ToEntity())

	g.Status(200)
}

// @Summary delete todo
// @Schemes
// @Description delete todo
// @Tags example
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Router /todos/{id} [delete]
// @Param id path string true "todo id"
func (h *Handler) DeleteTodo(g *gin.Context) {
	id := g.Param("id")

	err := h.TodoRepository.DeleteTodo(g.Request.Context(), id)

	if err != nil {
		g.Status(400)
		return
	}
	g.Status(200)

}

// @Summary list todo
// @Schemes
// @Description list todo
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {array} model.TodoItem
// @Failure 400
// @Router /todos [get]
func (h *Handler) ListTodos(g *gin.Context) {
	claims := GetClaims(g)

	log.Info().
		Strs("scopes", claims.AllScope).
		Str("claims", claims.Scope).
		Msg("i received claims")

	data, err := h.TodoRepository.ListTodo(g.Request.Context())

	if err != nil {
		g.Status(400)
		return
	}

	res := model.ToTodoItemSlice(data)

	log.Info().Any("data", res).Msg("gin response")

	g.JSON(200,
		res,
	)
}

// @Summary show detail
// @Schemes
// @Description show detail
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} model.TodoItem
// @Failure 400
// @Failure 404
// @Router /todos/{id} [get]
// @Param id path string true "todo id"
func (h *Handler) TodoDetail(g *gin.Context) {
	id := g.Param("id")
	item, err := h.TodoRepository.GetDetail(g.Request.Context(), id)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			g.Status(404)
			return
		}
		g.JSON(400, err.Error())
		return
	}

	g.JSON(200, item)

}

func GetClaims(g *gin.Context) middleware.AuthClaims {
	c := g.MustGet("claims")
	if fc, ok := c.(middleware.AuthClaims); ok {
		return fc
	}

	panic("no claims")

}
