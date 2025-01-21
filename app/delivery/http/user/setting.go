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

type settingHandler struct {
	SettingUsecase usecase_user.SettingUsecase
	Route          *gin.RouterGroup
	Middleware     middleware.AuthMiddleware
}

func NewSettingHandler(ginEngine *gin.Engine, middleware middleware.AuthMiddleware, settingUsecase usecase_user.SettingUsecase) {
	handler := &settingHandler{
		SettingUsecase: settingUsecase,
		Route:          ginEngine.Group("/user"),
		Middleware:     middleware,
	}

	handler.handleSettingRoute("/setting")
}

func (h *settingHandler) handleSettingRoute(path string) {
	api := h.Route.Group(path)

	api.PUT("/update-profile", h.Middleware.AuthUser(), h.UpdateProfile)
}

func (r *settingHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(model.JWTClaimUser)
	payload := request.UserUpdateProfileRequest{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{
			Data:    nil,
			Message: "invalid json data",
			Status:  http.StatusBadRequest,
		})
		return
	}

	response := r.SettingUsecase.UpdateProfile(ctx, claim, payload)

	c.JSON(response.Status, response)
}
