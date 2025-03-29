package auth

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=16"`
}

type RefreshRequest = tokenIncludedRequest

type LogoutRequest = tokenIncludedRequest

type tokenIncludedRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`	
}
