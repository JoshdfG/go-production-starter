package http

import (
	"net/http"

	"todo-clean/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	uc *usecase.TodoUseCase
}

func NewTodoHandler(uc *usecase.TodoUseCase) *TodoHandler {
	return &TodoHandler{uc: uc}
}

func (h *TodoHandler) Register(r *gin.RouterGroup) {
	todos := r.Group("/todos")
	{
		todos.GET("", h.list)
		todos.POST("", h.create)
		todos.PATCH("/:id", h.complete)
		todos.DELETE("/:id", h.delete)
	}
}

// @Summary     List all todos
// @Description Get all todos ordered by creation date
// @Tags        todos
// @Produce     json
// @Success     200 {array}  entity.Todo
// @Failure     500 {object} ErrorResponse
// @Router      /todos [get]
// @Security BearerAuth
// @Router /todos [get]
func (h *TodoHandler) list(c *gin.Context) {
	todos, err := h.uc.ListTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

// @Summary     Create a todo
// @Description Create a new todo item
// @Tags        todos
// @Accept      json
// @Produce     json
// @Param       request body     CreateTodoRequest true "Todo title"
// @Success     201     {object} entity.Todo
// @Failure     400     {object} ErrorResponse
// @Router      /todos [post]
// @Security BearerAuth
// @Router /todos [post]
func (h *TodoHandler) create(c *gin.Context) {
	var body CreateTodoRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	todo, err := h.uc.CreateTodo(body.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

// @Summary     Complete a todo
// @Description Mark a todo as done
// @Tags        todos
// @Produce     json
// @Param       id  path string true "Todo ID"
// @Success     204
// @Failure     404 {object} ErrorResponse
// @Router      /todos/{id} [patch]
// @Security BearerAuth
// @Router /todos/{id} [patch]
func (h *TodoHandler) complete(c *gin.Context) {
	id := c.Param("id")
	if err := h.uc.CompleteTodo(id); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary     Delete a todo
// @Description Permanently delete a todo
// @Tags        todos
// @Param       id  path string true "Todo ID"
// @Success     204
// @Failure     404 {object} ErrorResponse
// @Router      /todos/{id} [delete]
// @Security BearerAuth
// @Router /todos/{id} [delete]
func (h *TodoHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.uc.DeleteTodo(id); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
