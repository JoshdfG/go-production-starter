package http

// ErrorResponse is returned on all error responses
type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

// CreateTodoRequest is the request body for creating a todo
type CreateTodoRequest struct {
	Title string `json:"title" binding:"required" example:"Buy groceries"`
}

// AuthRequest is used for both register and login
type AuthRequest struct {
	Email    string `json:"email"    binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8"  example:"securepassword"`
}

// AuthResponse contains the JWT token
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGci..."`
}
