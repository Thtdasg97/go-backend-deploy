package usecase

import (
	"context"
	"go-backend/internal/dto"
)

type ChatMessageUsecase interface {
	FindAll(ctx context.Context, input dto.ChatMessageFindAllInput) (any, error)
}
