package handler

import (
	"context"
	"fmt"
	"go-backend/internal/dto"
	"go-backend/internal/usecase"
	"time"

	server "github.com/zishang520/socket.io/servers/socket/v3"
)

type ChatHandler struct {
	chatUsecase usecase.ChatUsecase
}

func NewChatHandler(chatUsecase usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{
		chatUsecase: chatUsecase,
	}
}

func (c *ChatHandler) CreateGroup(args ...any) {
	fmt.Printf("received: %v\n", args)
	fmt.Printf("received: %T | %v\n", args[0], args[0])

	payload := args[0].(map[string]interface{})

	accessToken := payload["accessToken"].(string)
	targetUserIdsAny := payload["targetUserIds"].([]interface{})
	ack := args[1].(func([]interface{}, error))

	name := ""
	nameAny := payload["name"]
	if nameAny != nil {
		name = nameAny.(string)
	}

	targetUserIds := []int{}
	for _, userIdAny := range targetUserIdsAny {
		fmt.Printf("received: %T | %v\n", userIdAny, userIdAny)
		userId := int(userIdAny.(float64))
		targetUserIds = append(targetUserIds, userId)

	}

	fmt.Println("accessToken", accessToken)
	fmt.Println("targetUserIds", targetUserIds)

	chatGroup, err := c.chatUsecase.CreateGroup(context.Background(), accessToken, targetUserIds, name)
	if err != nil {
		res := dto.ChatRes{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		ack([]interface{}{res}, nil)
		return

	}
	res := dto.ChatRes{
		Status:  "success",
		Message: "Create Chat Group Success",
		Data: map[string]any{
			"chatGroupId": chatGroup.ID,
		},
	}
	ack([]interface{}{res}, nil)

}

// luôn bắt gửi access token
func (c *ChatHandler) JoinGroup(socket *server.Socket, args ...any) {
	// gửi payload
	payload := args[0].(map[string]interface{})

	accessToken := payload["accessToken"].(string)
	chatGroupId := int(payload["chatGroupId"].(float64))
	ack := args[1].(func([]interface{}, error))

	// truyền 2 thằng accessToken vs chatGroupID vào usecase
	fmt.Println("accessToken", accessToken)
	fmt.Println("chatGroupId", chatGroupId)
	chatGroupExists, err := c.chatUsecase.JoinGroup(context.Background(), accessToken, chatGroupId)

	if err != nil {
		res := dto.ChatRes{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		ack([]interface{}{res}, nil)
		return

	}

	socket.Join(createNameRoom(chatGroupExists.ID))
	println("socket.Rooms()", socket.Rooms()) // xem coi có join được ko

	res := dto.ChatRes{
		Status:  "success",
		Message: "Join Group Success",
		Data:    nil,
	}
	ack([]interface{}{res}, nil)

	// fmt.Println("ack", ack)

}

// thiết kế nhận thêm tham số io, vì có io mới truyền vào đúng cái room được
func (c *ChatHandler) SendMessage(io *server.Server, args ...any) {

	payload := args[0].(map[string]interface{})
	accessToken := payload["accessToken"].(string)
	chatGroupId := int(payload["chatGroupId"].(float64))
	message := payload["message"].(string)
	ack := args[1].(func([]interface{}, error))

	fmt.Println("accesstoken", accessToken)
	fmt.Println("chatGroupId", chatGroupId)

	createdAt := time.Now().UTC()

	result, err := c.chatUsecase.SendMessage(context.Background(), accessToken, chatGroupId, message, createdAt)

	if err != nil {
		res := dto.ChatRes{
			Status:  "success",
			Message: "Join Group Success",
			Data:    nil,
		}
		ack([]interface{}{res}, nil)
		return
	}

	io.To(createNameRoom(result.ChatGroupId)).Emit("SEND_MESSAGE", map[string]any{
		"messageText": result.MessageText,
		"userId":      result.UserId,
		"chatGroupId": result.ChatGroupId,
		"createdAt":   createdAt.Format(time.RFC3339),
	})

}

func createNameRoom(chatGroupId int) server.Room {
	return server.Room(fmt.Sprintf("chat:%d", chatGroupId))
}
