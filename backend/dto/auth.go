package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type MeResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
