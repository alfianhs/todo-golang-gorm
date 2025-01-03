package http_user

import (
	"golang-gorm/app/delivery/http/middleware"
	usecase_user "golang-gorm/app/usecase/user"
	"golang-gorm/domain/model"
	"golang-gorm/domain/request"
	"golang-gorm/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type todoHandler struct {
	TodoUsecase usecase_user.TodoUsecase
	Route       *gin.RouterGroup
	Middleware  middleware.AuthMiddleware
}

func NewTodoHandler(ginEngine *gin.Engine, middleware middleware.AuthMiddleware, todoUsecase usecase_user.TodoUsecase) {
	handler := &todoHandler{
		TodoUsecase: todoUsecase,
		Route:       ginEngine.Group("/user"),
		Middleware:  middleware,
	}

	handler.handleTodoRoute("/todo")
}

func (h *todoHandler) handleTodoRoute(path string) {
	api := h.Route.Group(path)

	api.GET("", h.Middleware.AuthUser(), h.List)
	api.GET("/:id", h.Middleware.AuthUser(), h.GetByID)
	api.POST("", h.Middleware.AuthUser(), h.Create)
	api.PUT("/:id", h.Middleware.AuthUser(), h.Update)
	api.DELETE("/:id", h.Middleware.AuthUser(), h.Delete)
}

func (r *todoHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(model.JWTClaimUser)
	query := c.Request.URL.Query()

	response := r.TodoUsecase.GetAll(ctx, claim, query)

	c.JSON(response.Status, response)
}

func (r *todoHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(model.JWTClaimUser)
	todoID := c.Param("id")

	response := r.TodoUsecase.GetOne(ctx, claim, todoID)

	c.JSON(response.Status, response)
}

func (r *todoHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(model.JWTClaimUser)
	payload := request.CreateTodoRequest{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{
			Data:    nil,
			Message: "invalid json data",
			Status:  http.StatusBadRequest,
		})
		return
	}

	response := r.TodoUsecase.Create(ctx, claim, payload)

	c.JSON(response.Status, response)
}

func (r *todoHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(model.JWTClaimUser)
	todoID := c.Param("id")
	payload := request.UpdateTodoRequest{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{
			Data:    nil,
			Message: "invalid json data",
			Status:  http.StatusBadRequest,
		})
		return
	}

	response := r.TodoUsecase.Update(ctx, claim, todoID, payload)

	c.JSON(response.Status, response)
}

func (r *todoHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(model.JWTClaimUser)
	todoID := c.Param("id")

	response := r.TodoUsecase.Delete(ctx, claim, todoID)

	c.JSON(response.Status, response)
}
