/*
@Author: Weny Xu
@Date: 2021/06/02 22:41
*/

package odm

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type IModel interface {
	PrepareID(id interface{}) (interface{}, error)

	GetID() interface{}
	SetID(id interface{})
	Deleting(ctx context.Context, cfg *FieldsConfig) error
	Creating(ctx context.Context, cfg *FieldsConfig) error
	Saving(ctx context.Context, cfg *FieldsConfig) error
}

type DefaultModel struct {
	IDField         `bson:",inline"`
	TimestampFields `bson:",inline"`
	DeletedAtField  `bson:",inline"`
}

// CollectionGetter interface contains a method to return
// a model's custom collection.
type CollectionGetter interface {
	// Collection method return collection
	Collection() *mongo.Collection
}

// CollectionNameGetter interface contains a method to return
// the collection name of a model.
type CollectionNameGetter interface {
	// CollectionName method return model collection's name.
	CollectionName() string
}
