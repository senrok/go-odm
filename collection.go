/*
@Author: Weny Xu
@Date: 2021/06/02 22:42
*/

package odm

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	InvalidatedModels        = errors.New("invalidated models input")
	defaultInsertManyOptions = options.InsertMany().SetOrdered(true)
)

type Collection struct {
	opts *Options
	*mongo.Collection
	metadata     string
	fieldsConfig *FieldsConfig
}

func (o *Options) loadColl(meta string) (bool, *Collection) {
	if v, ok := o.collections.Load(meta); ok {
		if coll, ok := v.(*Collection); ok {
			return true, coll
		}
	}
	return false, nil
}

func (o *Options) storeColl(meta string, collection *Collection) {
	o.collections.Store(meta, collection)
}

// NewCollection returns a new collection using the current configuration values.
func (o *Options) NewCollection(m IModel, name string, opts ...*options.CollectionOptions) *Collection {
	var coll *mongo.Collection
	if getter, ok := m.(CollectionGetter); ok {
		coll = getter.Collection()
	} else {
		coll = o.db.Collection(name, opts...)
	}
	return &Collection{
		Collection:   coll,
		metadata:     name,
		opts:         o,
		fieldsConfig: genModelFieldsInfo(m).genModelConfig(),
	}
}

// Coll gets a collection or return a new collection
func (o *Options) Coll(m IModel) *Collection {
	meta := CollName(m)
	// exists
	if ok, coll := o.loadColl(meta); ok {
		return coll
		// not exists
	} else {
		coll := o.NewCollection(m, meta)
		o.storeColl(meta, coll)
		return coll
	}
}

// Create
func (c *Collection) Create(model IModel, opts ...*options.InsertOneOptions) error {
	ctx, _ := context.WithTimeout(context.Background(), c.opts.timeout)
	res, err := c.Collection.InsertOne(ctx, model, opts...)

	if err != nil {
		return err
	}

	// Set new id
	model.SetID(res.InsertedID)
	return err
}

// CreateMany
func (c *Collection) CreateMany(input interface{}, opts ...*options.InsertManyOptions) error {
	ctx, _ := context.WithTimeout(context.Background(), c.opts.timeout)
	models, err := prepareModels(ctx, c.fieldsConfig, input)
	if err != nil {
		return err
	}

	// run creating saving hooks
	if err = modelsHooksRunnerExecutor(ctx, c.fieldsConfig, models, creatingHook, savingHook); err != nil {
		return err
	}

	// set order to true forcefully
	result, err := c.Collection.InsertMany(ctx, models.Interfaces(), append(opts, defaultInsertManyOptions)...)
	if err != nil {
		return err
	}
	for index, id := range result.InsertedIDs {
		models.At(index).SetID(id)
	}
	return nil
}
