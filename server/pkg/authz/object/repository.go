package authz

import (
	"context"
	"fmt"

	"github.com/warrant-dev/warrant/server/pkg/database"
	"github.com/warrant-dev/warrant/server/pkg/middleware"
	"github.com/warrant-dev/warrant/server/pkg/service"
)

type ObjectRepository interface {
	Create(ctx context.Context, object Object) (int64, error)
	GetById(ctx context.Context, id int64) (*Object, error)
	GetByObjectTypeAndId(ctx context.Context, objectType string, objectId string) (*Object, error)
	List(ctx context.Context, filterOptions *FilterOptions, listParams middleware.ListParams) ([]Object, error)
	DeleteByObjectTypeAndId(ctx context.Context, objectType string, objectId string) error
}

func NewRepository(db database.Database) (ObjectRepository, error) {
	switch db.Type() {
	case database.TypeMySQL:
		mysql, ok := db.(*database.MySQL)
		if !ok {
			return nil, service.NewInternalError("Invalid database provided")
		}

		return NewMySQLRepository(mysql), nil
	default:
		return nil, service.NewInternalError(fmt.Sprintf("Invalid database type %s specified", db.Type()))
	}
}