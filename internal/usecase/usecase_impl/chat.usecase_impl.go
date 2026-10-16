package usecase_impl

import (
	"context"
	"errors"
	"fmt"
	"go-backend/ent"
	"go-backend/internal/dto"
	"go-backend/internal/repository"
	"go-backend/internal/usecase"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ChatUsecase struct {
	tokenUsecase              usecase.TokenUsecase
	userRepository            repository.UserRepository
	chatGroupRepository       repository.ChatGroupRepository
	chatGroupMemberRepository repository.ChatGroupMemberRepository
	unitOfWorkRepository      repository.UnitOfWorkRepository
	chatMessageRepository     repository.ChatMessageRepository
}

func NewChatUsecase(tokenUsecase usecase.TokenUsecase,
	userRepository repository.UserRepository,
	chatGroupRepository repository.ChatGroupRepository,
	chatGroupMemberRepository repository.ChatGroupMemberRepository,
	unitOfWorkRepository repository.UnitOfWorkRepository,
	chatMessageRepository repository.ChatMessageRepository) usecase.ChatUsecase {
	return &ChatUsecase{
		tokenUsecase:              tokenUsecase,
		userRepository:            userRepository,
		chatGroupRepository:       chatGroupRepository,
		chatGroupMemberRepository: chatGroupMemberRepository,
		unitOfWorkRepository:      unitOfWorkRepository,
		chatMessageRepository:     chatMessageRepository,
	}
}

// CreateGroup implements [usecase.ChatUsecase].
func (c *ChatUsecase) CreateGroup(ctx context.Context, accessToken string, targetUserIds []int, name string) (*ent.ChatGroups, error) {
	fmt.Println("accessToken", accessToken)

	claimAccessToken, err := c.tokenUsecase.VerifyAccessToken(accessToken, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}
	fmt.Println("accessToken", claimAccessToken)
	user, err := c.userRepository.FindUserById(ctx, claimAccessToken.UserId)

	if err != nil {
		return nil, err
	}

	fmt.Println("user", user)

	targetUserIds = append(targetUserIds, user.ID)

	fmt.Println("targetUserIds", targetUserIds)

	targetUserIdsUnique := []int{}

	seen := map[int]bool{}

	for _, id := range targetUserIds {
		if seen[id] {
			continue
		}

		seen[id] = true
		targetUserIdsUnique = append(targetUserIdsUnique, id)
	}
	fmt.Println("targetUserIdsUnique", targetUserIdsUnique)

	// xử lý logic để phân biệt đang chat nhiều hay 1-1

	var chatGroupExists *ent.ChatGroups

	if len(targetUserIdsUnique) == 2 {
		// chat 1 - 1
		chatGroupExists, err = c.chatGroupRepository.CheckChatGroupOneOneExist(ctx, targetUserIdsUnique)

		if err != nil && !ent.IsNotFound(err) {
			return nil, err
		}

		if chatGroupExists == nil {
			// unit of work = transaction
			err := c.unitOfWorkRepository.Do(ctx, func(ctxTx context.Context) error {
				chatGroupExists, err = c.chatGroupRepository.CreateGroup(ctxTx, "", user.ID)
				if err != nil {
					return err
				}

				err = c.chatGroupMemberRepository.CreateGroupMemberMany(ctxTx,
					chatGroupExists.ID,
					targetUserIdsUnique)

				if err != nil {
					return err
				}
				return nil
			})

			if err != nil {
				return nil, err
			}

		}

	} else {
		// chat group
		err := c.unitOfWorkRepository.Do(ctx, func(ctxTx context.Context) error {
			chatGroupExists, err = c.chatGroupRepository.CreateGroup(ctxTx, name, user.ID)
			if err != nil {
				return err
			}

			err = c.chatGroupMemberRepository.CreateGroupMemberMany(ctxTx,
				chatGroupExists.ID,
				targetUserIdsUnique)

			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

	}

	return chatGroupExists, nil
}

// JoinGroup implements [usecase.ChatUsecase].
func (c *ChatUsecase) JoinGroup(ctx context.Context, accessToken string, chatGroupId int) (*ent.ChatGroups, error) {
	claimAccessToken, err := c.tokenUsecase.VerifyAccessToken(accessToken, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}
	fmt.Println("accessToken", claimAccessToken)
	user, err := c.userRepository.FindUserById(ctx, claimAccessToken.UserId)

	if err != nil {
		return nil, err
	}

	// xử lý logic chatgroupID: tìm xem chatgroupID này có hay không và user có đang tồn tại trong nhóm chat đó ko

	chatGroup, err := c.chatGroupRepository.FindOneById(ctx, chatGroupId)
	if err != nil {
		return nil, err
	}

	isUserInChatGroup := false
	for _, member := range chatGroup.Edges.ChatGroupMembers {

		if member.Edges.Users.ID == user.ID {
			isUserInChatGroup = true
		}

	}

	if !isUserInChatGroup {
		return nil, errors.New("User not exists chat group")
	}

	return chatGroup, nil
}

// SendMessage implements [usecase.ChatUsecase].
func (c *ChatUsecase) SendMessage(ctx context.Context, accessToken string, chatGroupId int, message string, createdAt time.Time) (*dto.SendMessageReturn, error) {

	claimAccessToken, err := c.tokenUsecase.VerifyAccessToken(accessToken, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}
fmt.Println("accessToken", claimAccessToken)
	user, err := c.userRepository.FindUserById(ctx, claimAccessToken.UserId)

	if err != nil {
		return nil, err
	}

	// xử lý logic chatgroupID: tìm xem chatgroupID này có hay không và user có đang tồn tại trong nhóm chat đó ko

	chatGroup, err := c.chatGroupRepository.FindOneById(ctx, chatGroupId)
	if err != nil {
		return nil, err
	}

	isUserInChatGroup := false
	for _, member := range chatGroup.Edges.ChatGroupMembers {

		if member.Edges.Users.ID == user.ID {
			isUserInChatGroup = true
		}

	}

	if !isUserInChatGroup {
		return nil, errors.New("User not exists chat group")
	}

	go func() {
		_, err = c.chatMessageRepository.CreateMessage(ctx, user.ID, chatGroup.ID, message, createdAt)

		if err != nil {
			return
		}
	}()

	return &dto.SendMessageReturn{
		MessageText: message,
		ChatGroupId: chatGroup.ID,
		UserId:      user.ID,
	}, nil

}
