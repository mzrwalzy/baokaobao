package model

type UserResponse struct {
	ID        int64  `json:"id"`
	OpenID    string `json:"openid"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Phone     string `json:"phone"`
}

type AdminResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

type LoginByWechatRequest struct {
	Code string `json:"code" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token string        `json:"token"`
	User  AdminResponse `json:"user"`
}

type DecryptPhoneRequest struct {
	Code string `json:"code" binding:"required"`
}

type DecryptPhoneResponse struct {
	Phone string `json:"phone"`
}
