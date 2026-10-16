package response

import (
	"net/http"
)

// tạo struct
// implement interface
// Những instance được new từ Exception sẽ là lỗi

type Exception struct {
	StatusCode int
	Message    string
}

// khi new Exception struct lên, thì sẽ có sawx Method Error()
func (e *Exception) Error() string {
	return e.Message
}

func NewBadRequestException(message ...string) *Exception {

	statusCode := http.StatusBadRequest
	messageDefault := http.StatusText(statusCode)

	if len(message) > 0 && message[0] != "" {
		messageDefault = message[0]
	}

	return &Exception{
		StatusCode: statusCode,
		Message:    messageDefault,
	}
}

func NewInternalServerErrorException(message ...string) *Exception {

	statusCode := http.StatusInternalServerError
	messageDefault := http.StatusText(statusCode)

	if len(message) > 0 && message[0] != "" {
		messageDefault = message[0]
	}

	return &Exception{
		StatusCode: statusCode,
		Message:    messageDefault,
	}
}

// not found
func NewNotFoundErrorException(message ...string) *Exception {

	statusCode := http.StatusNotFound
	messageDefault := http.StatusText(statusCode)

	if len(message) > 0 && message[0] != "" {
		messageDefault = message[0]
	}

	return &Exception{
		StatusCode: statusCode,
		Message:    messageDefault,
	}
}

// forbidden 403 (refresh token)
func NewForbiddenErrorException(message ...string) *Exception {

	statusCode := http.StatusForbidden
	messageDefault := http.StatusText(statusCode)

	if len(message) > 0 && message[0] != "" {
		messageDefault = message[0]
	}

	return &Exception{
		StatusCode: statusCode,
		Message:    messageDefault,
	}
}

// 401: logout người dùng
func NewUnauthorizedErrorException(message ...string) *Exception {

	statusCode := http.StatusUnauthorized
	messageDefault := http.StatusText(statusCode)

	if len(message) > 0 && message[0] != "" {
		messageDefault = message[0]
	}

	return &Exception{
		StatusCode: statusCode,
		Message:    messageDefault,
	}
}









