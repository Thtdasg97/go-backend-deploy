package swagger

import (
	"go-backend/ent"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

func article(reflector *openapi3.Reflector) error {

	// cách dùng này gọi là slope code? mình gọi tên như vậy có đúng không vậy?
	// cái này dùng trong trường hợp nào? phù hợp với nghiệp vụ thế nào? Hình như cái này trong Java không có thì phải?
	{
		// add api cho swagger? không biết hiểu đúng không, sau khi code phần này xong thì F5 lại sẽ hiện list API,
		// mà không biết làm vậy để làm gì :))
		operationContext, err := reflector.NewOperationContext(http.MethodGet, "/api/article")
		if err != nil {
			return err
		}

		// Setting cái này để bên mà mình cung cấp không cần phải liên lạc lại gì với mình nữa, cung cấp thông tin đầy đủ
		// mục tiêu là như vậy. Đây là doc, trình bày đầy đủ (Summary, Description) mục đích là mô tả rõ API
		operationContext.SetTags("Article")
		operationContext.SetSummary("Lấy danh sách article")
		operationContext.SetDescription("có phân trang")

		operationContext.AddReqStructure(new(dto.ArticleFindallReq))
		operationContext.AddRespStructure(new(response.SuccessFormat[pagination.PaginationRes[*ent.Articles]]))
		reflector.AddOperation(operationContext)
	}

	{

		operationContext, err := reflector.NewOperationContext(http.MethodGet, "/api/article/{id}")
		if err != nil {
			return err
		}

		operationContext.SetTags("Article")
		operationContext.SetSummary("Lấy chi tiết article")
		operationContext.SetDescription("viết Description vào....")

		operationContext.AddReqStructure(new(struct {
			Id string `path:"id" description:"id của bài viết" example:"1"`
		}))
		operationContext.AddRespStructure(new(response.SuccessFormat[pagination.PaginationRes[*ent.Articles]]))
		reflector.AddOperation(operationContext)
	}
	return nil
}
