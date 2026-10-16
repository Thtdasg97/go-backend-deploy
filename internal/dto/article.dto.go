package dto

import (
	"go-backend/internal/common/pagination"
)

type ArticleCreateReq struct {
	Title string `json:"title" binding:"required"`
	// để * để có thể handle 2 giá trị, 1 là string, 2 là nil
	Content  *string `json:"content" binding:"omitempty"`
	ImageUrl *string `json:"image_url" binding:"omitempty"`
}

// bắt FE gửi theo format mình muốn thì design ở struct

type ArticleFindAllFilters struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Views   *int   `json:"views"`
}

type ArticleFindAllInput struct {
	// Query pagination.Query

	// kỹ thuật embedde struct

	// Page     int
	// PageSize int
	// Offset   int
	pagination.Query
	ArticleFindAllFilters
}

type ArticleFindallReq struct {
	pagination.Query
	Filters string `query:"filters" example:"{ \"id\":1, \"content\":\"string\", \"views\":55 }"`
}
