package repository

import "context"

type UnitOfWorkRepository interface {
	// nhét transaction vào context
	Do(ctx context.Context, fn func(ctxTx context.Context) error) error

	// callback function

}
