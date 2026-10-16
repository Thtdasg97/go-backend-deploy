package repository_impl

import (
	"context"
	"go-backend/internal/common/elastic"
	"go-backend/internal/repository"

	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/operator"
)

type searchRepository struct {
	elastic *elastic.Elastic
}

func NewSearchRepository(elastic *elastic.Elastic) repository.SearchRepository {
	return &searchRepository{
		elastic: elastic,
	}
}

// FindAll implements [repository.SearchRepository].
func (a *searchRepository) FindAll(ctx context.Context, textSearch string) (any, error) {

	result, err := a.elastic.EsClient.
		Search().
		Index("articles,users").
		Request(&search.Request{
			Query: &types.Query{
				// match nhiều loại
				MultiMatch: &types.MultiMatchQuery{
					Query: textSearch, // user gửi gì lên thì truyền vô đây
					Fields: []string{
						"title",
						"content",
						"email",
						"fullName",
					},
					// với nhiều từ khóa, chỉ cần khớp một phần là OK
					Operator: &operator.Or,

					// cho phép user gõ sai nhẹ vẫn tìm được
					Fuzziness: "AUTO",

					// Doucument nên khớp khoảng 60% số từ
					MinimumShouldMatch: "60%",

					// mấy thằng trên là tính năng xịn nhất của elastic, được xài nhiều nhất.

				},
			}}).Do(ctx)

	return result, err
}
