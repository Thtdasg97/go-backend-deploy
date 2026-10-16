package handler

import (
	"encoding/json"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ChatMessageHandler struct {
	chatMessageUsecase usecase.ChatMessageUsecase
}

func NewChatMessageHandler(chatMessageUsecase usecase.ChatMessageUsecase) *ChatMessageHandler {
	return &ChatMessageHandler{
		chatMessageUsecase: chatMessageUsecase,
	}
}

func (a *ChatMessageHandler) FindAll(ctx *gin.Context) {
	queryPagi := pagination.Get(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)

	// bắt FE gửi lên JSON
	filterString := ctx.DefaultQuery("filters", "{}")
	var filters dto.ChatMessageFindAllFilters
	json.Unmarshal([]byte(filterString), &filters)

	input := dto.ChatMessageFindAllInput{
		Query:                     queryPagi,
		ChatMessageFindAllFilters: filters,
	}

	result, err := a.chatMessageUsecase.FindAll(ctx.Request.Context(), input)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}
