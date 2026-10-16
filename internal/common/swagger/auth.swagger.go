package swagger

import (
	"go-backend/ent"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

func auth(reflector *openapi3.Reflector) error {
	{

		operationContext, err := reflector.NewOperationContext(http.MethodPost, "/api/auth/login")
		if err != nil {
			return err
		}

		operationContext.SetTags("Auth")
		operationContext.SetSummary("Lấy chi tiết article")
		operationContext.SetDescription("viết Description vào....")

		operationContext.AddReqStructure(new(dto.AuthLoginReq))
		operationContext.AddRespStructure(new(response.SuccessFormat[bool]))
		reflector.AddOperation(operationContext)
	}

	// api get-info
	{

		operationContext, err := reflector.NewOperationContext(http.MethodGet, "/api/auth/get-info")
		if err != nil {
			return err
		}

		operationContext.SetTags("Auth")
		operationContext.SetSummary("Lấy thông tin user")
		operationContext.SetDescription("viết description vào")

		operationContext.AddReqStructure(new(dto.AuthLoginReq))
		operationContext.AddRespStructure(new(response.SuccessFormat[*ent.Users]))
		reflector.AddOperation(operationContext)
	}

	// api refresh token
	{

		operationContext, err := reflector.NewOperationContext(http.MethodPost, "/api/auth/refresh-token")
		if err != nil {
			return err
		}

		operationContext.SetTags("Auth")
		operationContext.SetSummary("Lấy thông tin user")
		operationContext.SetDescription("viết description vào")

		operationContext.AddReqStructure(new(dto.AuthLoginReq))
		operationContext.AddRespStructure(new(response.SuccessFormat[bool]))
		reflector.AddOperation(operationContext)
	}

	return nil
}
