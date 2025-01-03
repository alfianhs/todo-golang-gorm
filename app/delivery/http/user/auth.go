package http_user

import (
	"net/http"

	"golang-gorm/app/delivery/http/middleware"
	usecase_user "golang-gorm/app/usecase/user"
	"golang-gorm/domain/model"
	"golang-gorm/domain/request"
	"golang-gorm/helpers"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	AuthUsecase usecase_user.AuthUsecase
	Route       *gin.RouterGroup
	Middleware  middleware.AuthMiddleware
}

func NewAuthHandler(ginEngine *gin.Engine, middleware middleware.AuthMiddleware, authUsecase usecase_user.AuthUsecase) {
	handler := &authHandler{
		AuthUsecase: authUsecase,
		Route:       ginEngine.Group("/user"),
		Middleware:  middleware,
	}

	handler.handleAuthRoute("/auth")
}

func (h *authHandler) handleAuthRoute(path string) {
	api := h.Route.Group(path)

	api.POST("/register", h.Register)
	api.POST("/login", h.Login)
	api.GET("/profile", h.Middleware.AuthUser(), h.GetProfile)
}

func (r *authHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.RegisterRequest{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	response := r.AuthUsecase.Register(ctx, payload)

	c.JSON(response.Status, response)
}

func (r *authHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.LoginRequest{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{
			Data:    nil,
			Message: "invalid json data",
			Status:  http.StatusBadRequest,
		})
		return
	}

	response := r.AuthUsecase.Login(ctx, payload)

	c.JSON(response.Status, response)
}

func (r *authHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(model.JWTClaimUser)

	response := r.AuthUsecase.GetProfile(ctx, claim)

	c.JSON(response.Status, response)
}
