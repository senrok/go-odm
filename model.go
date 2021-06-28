/*
@Author: Weny Xu
@Date: 2021/06/02 22:41
*/

package odm

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type IModel interface {
	PrepareID(id interface{}) (interface{}, error)

	GetID() interface{}
	SetID(id interface{})
	Deleting(ctx context.Context, cfg *FieldsConfig) error
	Creating(ctx context.Context, cfg *FieldsConfig) error
	Saving(ctx context.Context, cfg *FieldsConfig) error
}

type IModels []IModel

func (models IModels) Interfaces() (output []interface{}) {
	for v, _ := range models {
		output = append(output, v)
	}
	return
}

func (models IModels) At(i int) IModel {
	return models[i]
}

func prepareModels(ctx context.Context, cfg *FieldsConfig, models interface{}) (output IModels, err error) {
	rv := reflect.ValueOf(models)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			if model, ok := rv.Index(i).Interface().(IModel); ok {
				output = append(output, model)
			} else {
				return nil, InvalidatedModels
			}
		}
	default:
		return nil, InvalidatedModels
	}
	return output, nil
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
