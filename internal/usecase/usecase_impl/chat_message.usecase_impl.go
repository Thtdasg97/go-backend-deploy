package usecase_impl

import (
	"context"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/repository"
	"go-backend/internal/usecase"
	"math"
)

type chatMessageUsecase struct {
	chatMessageRepository repository.ChatMessageRepository
}

func NewChatMessageUsecase(chatMessageRepository repository.ChatMessageRepository) usecase.ChatMessageUsecase {
	return &chatMessageUsecase{
		chatMessageRepository: chatMessageRepository,
	}
}

// FindAll implements [usecase.ChatMessageUsecase].
func (a *chatMessageUsecase) FindAll(ctx context.Context, input dto.ChatMessageFindAllInput) (any, error) {
	data, err := a.chatMessageRepository.GetAll(ctx, input.Query, input.ChatMessageFindAllFilters)

	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	// totalItem: tổng số lượng item
	totalItem, err := a.chatMessageRepository.Count(ctx, input.ChatMessageFindAllFilters)
	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}
	// totalPage: tổng số trang = totalItem/ pageSize

	totalPage := float64(totalItem) / float64(input.PageSize)

	result := pagination.PaginationRes[any]{
		Items:     data,
		Page:      input.Page,
		PageSize:  input.PageSize,
		TotalItem: totalItem,
		TotalPage: int(math.Ceil(totalPage)),
	}

	return result, nil
}
