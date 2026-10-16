package handler

import (
	"encoding/json"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/usecase"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleUsecase usecase.ArticleUsecase
}

func NewArticleHandler(articleUsecase usecase.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{
		articleUsecase: articleUsecase,
	}
}

func (a *ArticleHandler) Create(ctx *gin.Context) {

	// design bắt FE gửi lên theo format mình đã code

	// nhận vào một pointer của 1 struct
	// mấy cái struct thì nhét vào DTO cho nó gọn

	// lấy body
	var body dto.ArticleCreateReq
	err := ctx.ShouldBindJSON(&body)

	if err != nil {

		if err == io.EOF {
			ctx.Error(response.NewBadRequestException("Body required homies"))
			return
		}

		ctx.Error(response.NewBadRequestException(err.Error()))
		return

	}

	result, err := a.articleUsecase.Create(ctx, body)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)

}

func (a *ArticleHandler) FindAll(ctx *gin.Context) {

	queryPagi := pagination.Get(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)

	// bắt FE gửi lên JSON
	filterString := ctx.DefaultQuery("filters", "{}")
	var filters dto.ArticleFindAllFilters
	json.Unmarshal([]byte(filterString), &filters)

	input := dto.ArticleFindAllInput{
		Query:                 queryPagi,
		ArticleFindAllFilters: filters,
	}

	result, err := a.articleUsecase.FindAll(ctx.Request.Context(), input)

	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)

}

func (a *ArticleHandler) Delete(ctx *gin.Context) {

	// get id > parse qua int vì trong database id kiểu int
	idString := ctx.Param("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	result, err := a.articleUsecase.Delete(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}
