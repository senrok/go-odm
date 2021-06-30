/*
@Author: Weny Xu
@Date: 2021/06/02 22:42
*/

package odm

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	InvalidatedModels        = errors.New("invalidated models input\n try []*YourModel")
	defaultInsertManyOptions = options.InsertMany().SetOrdered(true)
	UnableSoftDeletable      = errors.New("unable soft-delete")
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
func (c *Collection) Create(ctx context.Context, model IModel, opts ...*options.InsertOneOptions) error {
	if err := modelHooksRunnerExecutor(ctx, c.fieldsConfig, model, creatingHook, savingHook); err != nil {
		return err
	}

	res, err := c.Collection.InsertOne(ctx, model, opts...)

	if err != nil {
		return err
	}

	// Set new id
	model.SetID(res.InsertedID)

	if err = modelHooksRunnerExecutor(ctx, c.fieldsConfig, model, createdHook, savedHook); err != nil {
		return err
	}

	return err
}

// CreateMany
func (c *Collection) CreateMany(ctx context.Context, input interface{}, opts ...*options.InsertManyOptions) error {
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

	if err = modelsHooksRunnerExecutor(ctx, c.fieldsConfig, models, createdHook, savedHook); err != nil {
		return err
	}

	return nil
}

// UpdateOne
func (c *Collection) UpdateOne(ctx context.Context, model IModel, opts ...*options.UpdateOptions) error {
	filter := bson.M{c.fieldsConfig.PrimaryIDBsonField: model.GetID()}

	if c.fieldsConfig.SoftDeletable() {
		excludeSoftDeletedItems(c.fieldsConfig.DeleteTimeBsonField, filter)
	}

	if err := modelHooksRunnerExecutor(ctx, c.fieldsConfig, model, updatingHook, savingHook); err != nil {
		return err
	}

	res, err := c.Collection.UpdateOne(ctx, filter, bson.M{"$set": model}, opts...)
	if err != nil {
		return err
	}

	if err := modelHooksRunnerExecutor(ctx, c.fieldsConfig, model, updatedHook(res), savedHook); err != nil {
		return err
	}
	return nil
}

// Find
func (c *Collection) Find(ctx context.Context, filter bson.M, result interface{}, opts ...*options.FindOptions) error {
	if c.fieldsConfig.SoftDeletable() {
		excludeSoftDeletedItems(c.fieldsConfig.DeleteTimeBsonField, filter)
	}
	res, err := c.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	if err = res.All(ctx, result); err != nil {
		return err
	}
	return nil
}

// FindOne
func (c *Collection) FindOne(ctx context.Context, filter bson.M, result interface{}, opts ...*options.FindOneOptions) error {
	if c.fieldsConfig.SoftDeletable() {
		excludeSoftDeletedItems(c.fieldsConfig.DeleteTimeBsonField, filter)
	}
	return c.Collection.FindOne(ctx, filter, opts...).Decode(result)
}

// SoftDeleteOne
func (c *Collection) SoftDeleteOne(ctx context.Context, model IModel, opts ...*options.UpdateOptions) error {
	filter := bson.M{c.fieldsConfig.PrimaryIDBsonField: model.GetID()}

	if !c.fieldsConfig.SoftDeletable() {
		return UnableSoftDeletable
	} else {
		excludeSoftDeletedItems(c.fieldsConfig.DeleteTimeBsonField, filter)
	}

	if err := modelHooksRunnerExecutor(ctx, c.fieldsConfig, model, softDeletingHook); err != nil {
		return err
	}

	res, err := c.Collection.UpdateOne(ctx, filter, bson.M{"$set": model}, opts...)
	if err != nil {
		return err
	}

	if err := modelHooksRunnerExecutor(ctx, c.fieldsConfig, model, softDeletedHook(res)); err != nil {
		return err
	}
	return nil
}

// SoftDeleteMany
func (c *Collection) SoftDeleteMany(ctx context.Context, filter bson.M, opts ...*options.UpdateOptions) error {
	if c.fieldsConfig.SoftDeletable() {
		excludeSoftDeletedItems(c.fieldsConfig.DeleteTimeBsonField, filter)
	}
	_, err := c.Collection.UpdateMany(ctx, filter, bson.M{"$set": bson.M{c.fieldsConfig.DeleteTimeBsonField: time.Now()}})
	if err != nil {
		return err
	}
	return nil
}

// RestoreOne
func (c *Collection) RestoreOne(ctx context.Context, model IModel, opts ...*options.UpdateOptions) error {
	filter := bson.M{c.fieldsConfig.PrimaryIDBsonField: model.GetID()}

	if !c.fieldsConfig.SoftDeletable() {
		return UnableSoftDeletable
	} else {
		onlySoftDeletedItems(c.fieldsConfig.DeleteTimeBsonField, filter)
	}

	if err := modelHooksRunnerExecutor(ctx, c.fieldsConfig, model, restoringHook); err != nil {
		return err
	}

	res, err := c.Collection.UpdateOne(ctx, filter, bson.M{"$set": model}, opts...)
	if err != nil {
		return err
	}

	if err := modelHooksRunnerExecutor(ctx, c.fieldsConfig, model, restoredHook(res)); err != nil {
		return err
	}
	return nil
}

// RestoreMany
func (c *Collection) RestoreMany(ctx context.Context, filter bson.M, opts ...*options.UpdateOptions) error {
	if !c.fieldsConfig.SoftDeletable() {
		return UnableSoftDeletable
	} else {
		onlySoftDeletedItems(c.fieldsConfig.DeleteTimeBsonField, filter)
	}

	_, err := c.Collection.UpdateMany(ctx, filter, bson.M{"$unset": bson.M{c.fieldsConfig.DeleteTimeBsonField: ""}})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOne
func (c *Collection) DeleteOne(ctx context.Context, model IModel, opts ...*options.DeleteOptions) error {

	if err := modelHooksRunnerExecutor(ctx, c.fieldsConfig, model, deletingHook); err != nil {
		return err
	}

	filter := bson.M{c.fieldsConfig.PrimaryIDBsonField: model.GetID()}

	res, err := c.Collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return err
	}

	if err := modelHooksRunnerExecutor(ctx, c.fieldsConfig, model, deletedHook(res)); err != nil {
		return err
	}
	return nil
}

// DeleteMany
func (c *Collection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) error {
	_, err := c.Collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return err
	}
	return nil
}

// Count
func (c *Collection) Count(ctx context.Context, filter bson.M, opts ...*options.CountOptions) (int64, error) {
	if c.fieldsConfig.SoftDeletable() {
		excludeSoftDeletedItems(c.fieldsConfig.DeleteTimeBsonField, filter)
	}
	return c.Collection.CountDocuments(ctx, filter, opts...)
}

func excludeSoftDeletedItems(deletedAtField string, m bson.M) {
	if _, ok := m[deletedAtField]; ok {
	} else {
		m[deletedAtField] = bson.M{"$exists": false}
	}
}

func onlySoftDeletedItems(deletedAtField string, m bson.M) {
	m[deletedAtField] = bson.M{"$exists": true}
}
