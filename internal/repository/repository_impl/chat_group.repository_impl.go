package repository_impl

import (
	"context"
	"go-backend/ent"
	"go-backend/ent/chatgroupmembers"
	"go-backend/ent/chatgroups"
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
	"go-backend/internal/repository"
)

type chatGroupRepository struct {
	entClient *ent.Client
}

func NewchatGroupRepository(entClient *ent.Client) repository.ChatGroupRepository {
	return &chatGroupRepository{
		entClient: entClient,
	}
}

// FindAll implements [repository.chatGroupRepository].
func (a *chatGroupRepository) FindAll(ctx context.Context) (any, error) {
	return nil, nil
}

// CheckChatGroupOneOneExist implements [repository.ChatGroupRepository].
func (a *chatGroupRepository) CheckChatGroupOneOneExist(ctx context.Context, ids []int) (*ent.ChatGroups, error) {
	entClient := a.entClient.ChatGroups.Query()
	// mấy cái này chịu, chả hiểu mẹ gì hết :))
	entClient = entClient.Where(
		chatgroups.NameIsNil(),
		chatgroups.HasChatGroupMembersWith(
			chatgroupmembers.UserIDIn(ids...),
		),
	)
	return entClient.Only(ctx)
}

// CreateGroup implements [repository.ChatGroupRepository].
func (a *chatGroupRepository) CreateGroup(ctx context.Context, name string, userId int) (*ent.ChatGroups, error) {
	entClientTx := GetClientTx(ctx, a.entClient)

	entCreate := entClientTx.ChatGroups.Create()

	if name != "" {
		entCreate = entCreate.SetName(name)
	}

	entCreate = entCreate.SetUserID(userId)

	return entCreate.Save(ctx)
}

// GetAll implements [repository.ChatGroupRepository].
func (a *chatGroupRepository) GetAll(ctx context.Context, query pagination.Query, filters dto.ChatGroupFindAllFilters) (any, error) {
	entQuery := a.entClient.ChatGroups.Query()

	// đặt logic xử lý trường hợp user không gửi ID lên

	handlerFilterChatGroup(filters, entQuery)

	// include
	entQuery = entQuery.WithUsers()

	entQuery = entQuery.WithChatGroupMembers(func(cgmq *ent.ChatGroupMembersQuery) {
		cgmq.WithUsers()
	})

	entQuery = entQuery.Limit(query.PageSize)
	entQuery = entQuery.Offset(query.Offset)

	// sắp xếp group chat theo thứ tự giảm dần
	entQuery = entQuery.Order(ent.Desc(chatgroups.FieldCreatedAt))

	return entQuery.All(ctx)
}

// Count implements [repository.ChatGroupRepository].
func (a *chatGroupRepository) Count(ctx context.Context, filters dto.ChatGroupFindAllFilters) (int, error) {
	entQuery := a.entClient.ChatGroups.Query()

	// đặt logic xử lý trường hợp user không gửi ID lên
	handlerFilterChatGroup(filters, entQuery)

	return entQuery.Count(ctx)
}

// tạo hàm cục bọ để hạn chế lặp code
func handlerFilterChatGroup(filters dto.ChatGroupFindAllFilters, entQuery *ent.ChatGroupsQuery) {

	if filters.Name != "" {
		entQuery = entQuery.Where(chatgroups.NameContainsFold(filters.Name))
	}

}

// FindOneById implements [repository.ChatGroupRepository].
func (a *chatGroupRepository) FindOneById(ctx context.Context, id int) (*ent.ChatGroups, error) {

	entQuery := a.entClient.ChatGroups.Query()
	entQuery = entQuery.Where(chatgroups.IDEQ(id))

	// lấy thêm thông tin user trong chatgroup
	// tham số là một function thì có thuật ngữ nào nói về nó ko? 
	entQuery = entQuery.WithChatGroupMembers(func(cgmq *ent.ChatGroupMembersQuery) {
		cgmq.WithUsers()
	})

	return entQuery.Only(ctx)
}
