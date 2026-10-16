package handler

import (
	"errors"
	"go-backend/internal/common/constant"
	"go-backend/internal/common/env"
	"go-backend/internal/common/helpers"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/usecase"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	env         *env.Env
}

func NewAuthHandler(authUsecase usecase.AuthUsecase, env *env.Env) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		env:         env,
	}
}

// register
func (a *AuthHandler) Register(ctx *gin.Context) {
	// tạo DTO để gửi body
	// parse JSON body từ request vào struct
	// bắt lỗi body rỗng -> trả về body required
	var body dto.AuthRegisterReq
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		if err == io.EOF {

			ctx.Error(response.NewBadRequestException("Body required"))
			return
		}

		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	result, err := a.authUsecase.Register(ctx.Request.Context(), body)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}

// login logic
func (a *AuthHandler) Login(ctx *gin.Context) {

	// tạo DTO để gửi body
	// parse JSON body từ request vào struct
	// bắt lỗi body rỗng -> trả về body required
	var body dto.AuthLoginReq
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		if err == io.EOF {

			ctx.Error(response.NewBadRequestException("Body required"))
			return
		}

		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	result, err := a.authUsecase.Login(ctx.Request.Context(), body)
	if err != nil {
		ctx.Error(err)
		return
	}
	setTokenCookie(ctx, result.AccessToken, result.RefreshToken)

	response.Success(true, "", 0, ctx)
}

// get info logic
func (a *AuthHandler) GetInfo(ctx *gin.Context) {
	user, err := helpers.GetUser(ctx)
	if err != nil {
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	result, err := a.authUsecase.GetInfo(ctx.Request.Context(), user)
	if err != nil {

		ctx.Error(err)
		return
	}

	// làm cái Key user ở đây có ý đồ gì ta?
	response.Success(map[string]any{
		"user": result,
	}, "", 0, ctx)
}

// xử lý refresh token
func (a *AuthHandler) RefreshToken(ctx *gin.Context) {
	accessToken, err := ctx.Cookie("accessToken")

	if err != nil {
		// kiểm tra access token, nếu có vấn đề gì xảy ra -> Logout user đó luôn

		if errors.Is(err, http.ErrNoCookie) {
			ctx.Error(response.NewUnauthorizedErrorException())
			ctx.Abort()
			return
		}
		ctx.Error(response.NewBadRequestException())
		return
	}

	refreshToken, err := ctx.Cookie("refreshToken")

	if err != nil {

		if errors.Is(err, http.ErrNoCookie) {
			ctx.Error(response.NewUnauthorizedErrorException())
			ctx.Abort()
			return
		}
		ctx.Error(response.NewBadRequestException())
		return
	}

	result, err := a.authUsecase.RefreshToken(ctx.Request.Context(), accessToken, refreshToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	setTokenCookie(ctx, result.AccessToken, result.RefreshToken)

	response.Success(true, "", 0, ctx)
}

// xử lý login google
func (a *AuthHandler) GoogleLogin(ctx *gin.Context) {

	result, err := a.authUsecase.GoogleLogin(ctx.Request.Context())
	if err != nil {
		ctx.Redirect(http.StatusFound, a.env.DomainFe+"/Login?error="+err.Error())
		return
	}

	// set state cho người dùng
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		constant.GOOGLE_OAUTH_STATE,
		result.State,
		60,
		"/",
		"",
		false,
		true,
	)

	redirect := ctx.Query("redirect")
	if redirect != "" {
		ctx.SetCookie(
			constant.GOOGLE_OAUTH_REDIRECT,
			redirect,
			60,
			"/",
			"",
			false,
			true,
		)

	}

	// redirect người dùng
	ctx.Redirect(http.StatusFound, result.Url)

}

// xử lý callback
func (a *AuthHandler) GoogleCallback(ctx *gin.Context) {

	// XỬ LÝ STATE
	// query state trên URL
	// xử lý case khi google không trả về (làm chuẩn chỉnh thôi chứ google thường trả về)
	stateGoogle := ctx.Query("state")
	if stateGoogle == "" {
		// redirect user về trang login
		ctx.Redirect(http.StatusFound, a.env.DomainFe+"/login?error="+"state empty")
		return

	}
	// state cookie
	stateCookie, err := ctx.Cookie(constant.GOOGLE_OAUTH_STATE)
	if err != nil {
		ctx.Redirect(http.StatusFound, a.env.DomainFe+"/login?error="+err.Error())
		return
	}

	// XỬ LÝ CODE (trên URL)
	code := ctx.Query("code")
	if code == "" {
		ctx.Redirect(http.StatusFound, a.env.DomainFe+"/login?error="+"code empty")
		return
	}

	input := dto.AuthGoogleCallbackInput{

		StateGoogle: stateGoogle,
		StateCookie: stateCookie,
		Code:        code,
	}

	result, err := a.authUsecase.GoogleCallback(ctx.Request.Context(), input)
	if err != nil {
		ctx.Error(err)
		return
	}

	// set cookies -> redirect
	setTokenCookie(ctx, result.AccessToken, result.RefreshToken)

	// state cookie
	redirectUrl, _ := ctx.Cookie(constant.GOOGLE_OAUTH_REDIRECT)

	// khác rỗng = không lỗi, mà không lỗi thì nghĩa là có giá trị, mà có giá trị thì redirect về URL
	if redirectUrl != "" {
		ctx.Redirect(http.StatusFound, redirectUrl)
		return
	}

	ctx.Redirect(http.StatusFound, a.env.DomainFe)

	response.Success(result, "", 0, ctx)
}

func setTokenCookie(ctx *gin.Context, accessToken string, refreshToken string) {

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"accessToken",
		accessToken,
		0,
		"/",
		"",
		false,
		false,
	)

	ctx.SetCookie(
		"refreshToken",
		refreshToken,
		0,
		"/",
		"",
		false,
		false,
	)

}
