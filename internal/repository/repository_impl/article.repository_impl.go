package repository_impl

import (
	"context"
	"go-backend/ent"
	"go-backend/ent/articles"
	"go-backend/internal/common/models"
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
	"go-backend/internal/repository"

	"gorm.io/gorm"
)

type articleRepository struct {
	entClient  *ent.Client
	gormClient *gorm.DB
}

func NewArticleRepository(entClient *ent.Client, gormClient *gorm.DB) repository.ArticleRepository {
	return &articleRepository{
		entClient:  entClient,
		gormClient: gormClient,
	}
}

// Create implements [repository.ArticleRepository].
func (a *articleRepository) Create(ctx context.Context, body dto.ArticleCreateReq) (any, error) {
	// set giá trị vô database
	entCreate := a.entClient.Articles.Create()

	entCreate = entCreate.SetTitle(body.Title)

	if body.Content != nil {
		// * + pointer = set gias trị
		entCreate = entCreate.SetContent(*body.Content)
	}

	if body.ImageUrl != nil {
		// * + pointer = set gias trị
		entCreate = entCreate.SetContent(*body.ImageUrl)
	}

	entCreate = entCreate.SetUserID(1)

	return entCreate.Save(ctx)

}

// // Create implements [usecase.articleRepository].
// func (a *articleRepository) Create(ctx context.Context) (any, error) {
//
//		a.entClient.Articles.Create().SetTitle("").SetContent("").SetImageURL("").Save(ctx)
//		return "Create", nil
//	}
//
// CreateGorm implements [repository.ArticleRepository].
func (a *articleRepository) CreateGorm(ctx context.Context, body dto.ArticleCreateReq) (any, error) {

	article := models.Article{
		Title:    body.Title,
		Content:  body.Content,
		ImageUrl: body.ImageUrl,
	}

	result := gorm.WithResult()

	err := gorm.G[models.Article](a.gormClient, result).Create(ctx, &article)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAll implements [repository.ArticleRepository].
func (a *articleRepository) GetAll(ctx context.Context, query pagination.Query, filters dto.ArticleFindAllFilters) ([]*ent.Articles, error) {
	entQuery := a.entClient.Articles.Query()

	// đặt logic xử lý trường hợp user không gửi ID lên

	handlerFilter(filters, entQuery)

	// include
	entQuery = entQuery.WithUsers()

	entQuery = entQuery.Limit(query.PageSize)

	entQuery = entQuery.Offset(query.Offset)
	return entQuery.All(ctx)
}

// Count implements [repository.ArticleRepository].
// implement để số trang fit với filter (đúng chuẩn UI-UX)
func (a *articleRepository) Count(ctx context.Context, filters dto.ArticleFindAllFilters) (int, error) {
	entQuery := a.entClient.Articles.Query()

	// đặt logic xử lý trường hợp user không gửi ID lên
	handlerFilter(filters, entQuery)

	return entQuery.Count(ctx)
}

// tạo hàm cục bọ để hạn chế lặp code
func handlerFilter(filters dto.ArticleFindAllFilters, entQuery *ent.ArticlesQuery) {

	if filters.Id > 0 {

		entQuery = entQuery.Where(articles.IDEQ(filters.Id))
	}

	if filters.Content != "" {
		entQuery = entQuery.Where(articles.ContentContainsFold(filters.Content))
	}
	if filters.Views != nil {
		entQuery = entQuery.Where(articles.ViewsEQ(*filters.Views))
	}
}

// Delete implements [repository.ArticleRepository].
func (a *articleRepository) Delete(ctx context.Context, id int) error {
	entDelete := a.entClient.Articles.DeleteOneID(id)

	return entDelete.Exec(ctx)
}
