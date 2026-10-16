package handler

import (
	"fmt"
	"go-backend/internal/common/middlewares"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Resonpsibility: Nhận dữ liệu, validate data đầu vào, chuyển cho usecase xử lý logic

// truyền interface
type DemoHandler struct {
	demoUsecase usecase.DemoUsecase
}

// type DemoHandler struct {
// 	demoUsecase *usecase.DemoUsecase
// }

func NewDemoHandler(demoUsecase usecase.DemoUsecase) *DemoHandler {
	return &DemoHandler{
		demoUsecase: demoUsecase,
	}
}

// func này có thể được gọi là middleware
func (d *DemoHandler) Query(ctx *gin.Context) {
	// lấy data FE gửi

	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))

	// sử dụng struct query
	input := pagination.Query{
		Page:     page,
		PageSize: pageSize,
	}

	// tham số nhận vào
	data := d.demoUsecase.Query(input)
	response.Success(data, "", 0, ctx)
}

func (d *DemoHandler) Param(ctx *gin.Context) {

	// demo chuyền data giữa các middleware --- start
	// lấy data trong storage ra
	userRaw, exists := ctx.Get("user")

	if !exists {
		fmt.Println("Lỗi không có user")
	}

	// parse qua struct, vì ctx.get trả any

	user, ok := userRaw.(middlewares.User)
	if !ok {
		fmt.Println("Lỗi không phải struct User")
	}

	fmt.Println("Handler đã nhận được: ", user)
	// demo chuyền data giữa các middleware --- end

	// xử lý lỗi
	id, err := strconv.Atoi(ctx.Param("id"))

	// demo lỗi không kiểm soát
	// panic("Lỗi không kiểm soát")

	// error khác nil thì dùng .Error để push lỗi vào
	if err != nil {
		// thay vì nhét err vô đây, mình sẽ nhét lỗi do mình custom vô
		ctx.Error(response.NewBadRequestException(""))
		// nếu không ngắt dòng (return) ở đây sẽ bị lỗi hiển thị 2 object
		return
	}

	data := d.demoUsecase.Param(id)
	response.Success(data, "", 0, ctx)

}

// Body
func (d *DemoHandler) Body(ctx *gin.Context) {

	var body dto.DemoBody

	// hàm này cũng có validate, print lỗi ra
	err := ctx.ShouldBindJSON(&body)
	fmt.Println("err", err)

	data := d.demoUsecase.Body(body)
	response.Success(data, "", 0, ctx)

}

// Header
func (d *DemoHandler) Header(ctx *gin.Context) {

	apiKey := ctx.GetHeader("API_KEY")

	// Response
	data := d.demoUsecase.Header(apiKey)
	response.Success(data, "", 0, ctx)

}
