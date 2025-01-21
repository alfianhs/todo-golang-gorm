package request

type UserUpdateProfileRequest struct {
	Name           string `json:"name" validate:"required"`
	ProfilePicture string `json:"profile_picture" validate:"required"`
}
