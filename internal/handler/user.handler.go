package handler

import (
	"encoding/json"
	"go-backend/internal/common/helpers"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (a *UserHandler) FindAll(ctx *gin.Context) {
	queryPagi := pagination.Get(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)

	// bắt FE gửi lên JSON
	filterString := ctx.DefaultQuery("filters", "{}")
	var filters dto.UserFindAllFilters
	json.Unmarshal([]byte(filterString), &filters)

	input := dto.UserFindAllInput{
		Query:              queryPagi,
		UserFindAllFilters: filters,
	}
	result, err := a.userUsecase.FindAll(ctx.Request.Context(), input)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}

func (a *UserHandler) FindOne(ctx *gin.Context) {

	idString := ctx.Param("id")
	if idString == "" {
		ctx.Error(response.NewBadRequestException("Required id"))
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.Error(response.NewBadRequestException("id invalid"))
		return
	}

	result, err := a.userUsecase.FindOne(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}

func (a *UserHandler) AvatarLocal(ctx *gin.Context) {

	user, err := helpers.GetUser(ctx)

	if err != nil {
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	// ý đồ dùng FormFile ở đây là gì? chưa hiểu rõ lắm
	fileHeader, err := ctx.FormFile("avatar")

	if err != nil {
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	result, err := a.userUsecase.AvatarLocal(ctx.Request.Context(), fileHeader, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}

func (a *UserHandler) AvatarCloud(ctx *gin.Context) {

	user, err := helpers.GetUser(ctx)

	if err != nil {
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	// ý đồ dùng FormFile ở đây là gì? chưa hiểu rõ lắm
	fileHeader, err := ctx.FormFile("avatar")

	if err != nil {
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}
	result, err := a.userUsecase.AvatarCloud(ctx.Request.Context(), fileHeader, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}
