package dto

type DemoBody struct {
	// thêm struct tag
	// thêm validate
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}


