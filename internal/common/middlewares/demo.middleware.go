package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// tạo struct user demo cách chuyền dữ liệu giữa các middleware

type User struct {
	Id    int
	Email string
	Pass  string
}

func A(ctx *gin.Context) {
	fmt.Println("Mid A Before")
	// khởi tạo struct TRONG middleware A
	user := User{
		Id:    1,
		Email: "thanh@gmail.com",
		Pass:  "123",
	}
	// gắn data vào storage
	ctx.Set("user", user)
	ctx.Next()
	fmt.Println("Mid A After")

}

func B(ctx *gin.Context) {
	fmt.Println("Mid B Before")
	ctx.Next()
	fmt.Println("Mid B After")
}

func C(ctx *gin.Context) {
	fmt.Println("Mid C Before")

	// lấy data trong storage ra
	userRaw, exists := ctx.Get("user")

	if !exists {
		fmt.Println("Lỗi không có user")
	}

	// parse qua struct, vì ctx.get trả any

	user, ok := userRaw.(User)
	if !ok {
		fmt.Println("Lỗi không phải struct User")
	}

	fmt.Println("Mid C đã nhận được: ", user)

	ctx.Next()
	
	fmt.Println("Mid C After")
}
