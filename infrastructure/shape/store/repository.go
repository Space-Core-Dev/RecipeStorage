package store

import (
	"context"

	"github.com/Space-Core-Dev/Mistaken.Recipe/core/domain/entity"
	"github.com/Space-Core-Dev/Mistaken.Utils/result"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository[T any] interface {
	Get(ctx context.Context, coll *mongo.Collection, id entity.Id) result.Result[T]
	Update(ctx context.Context, coll *mongo.Collection, entity T) result.Result[result.Void]
	Add(ctx context.Context, coll *mongo.Collection, entity T) result.Result[result.Void]
	Delete(ctx context.Context, coll *mongo.Collection, id entity.Id) result.Result[result.Void]
}
