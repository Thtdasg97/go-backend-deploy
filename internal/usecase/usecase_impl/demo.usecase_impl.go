package usecase_impl

import (
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
	"go-backend/internal/usecase"
)

type DemoUsecase struct{}

// func NewDemoUsecase() *DemoUsecase {
// 	return &DemoUsecase{}
// }

// thay vì trả về con trỏ, thì trả về interface
// khi trả về phải đầy đủ phương phức đã qui định ở interface
func NewDemoUsecase() usecase.DemoUsecase {
	return &DemoUsecase{}
}

func (d *DemoUsecase) Query(input pagination.Query) any {
	return input
}

func (d *DemoUsecase) Param(id int) int {
	return 1
}

func (d *DemoUsecase) Body(body dto.DemoBody) any {
	return body
}

func (d *DemoUsecase) Header(apiKey string) any {
	return apiKey
}
