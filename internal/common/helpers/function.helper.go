package helpers

import (
	"errors"
	"go-backend/ent"

	"github.com/gin-gonic/gin"
)

func GetUser(ctx *gin.Context) (*ent.Users, error) {
	userAny, exists := ctx.Get("user")
	if !exists {
		return nil, errors.New("API chưa gắn middleware prottect")
	}

	if userAny == nil {
		return nil, errors.New("Không có user")
	}

	// đây ông thầy nói kỹ thuật gì ấy nghe ko rõ, type absort?
	user, ok := userAny.(*ent.Users)
	if !ok {
		return nil, errors.New("Type User không hợp leej")

	}
	return user, nil

}
