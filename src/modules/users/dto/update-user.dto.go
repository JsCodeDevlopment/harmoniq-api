package dto

type UpdateUserDto struct {
	Name string `json:"name"`
}

type ChangePasswordDto struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
