/*
@Author: Weny Xu
@Date: 2021/06/03 0:09
*/

package odm

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreatingHook is called before saving a new model to the database
type CreatingHook interface {
	Creating(ctx context.Context, cfg *FieldsConfig) error
}

// CreatedHook is called after a model has been created
type CreatedHook interface {
	Created(ctx context.Context) error
}

// UpdatingHook is called before updating a model
type UpdatingHook interface {
	Updating(ctx context.Context) error
}

// UpdatedHook is called after a model is updated
type UpdatedHook interface {
	Updated(ctx context.Context, result *mongo.UpdateResult) error
}

// SoftDeletingHook is called before soft deleting a model
type SoftDeletingHook interface {
	SoftDeleting(ctx context.Context) error
}

// SoftDeletedHook is called after soft deleting a model
type SoftDeletedHook interface {
	SoftDeleted(ctx context.Context, result *mongo.UpdateResult) error
}

// RestoringHook is called before soft restoring a model
type RestoringHook interface {
	Restoring(ctx context.Context) error
}

// RestoredHook is called after soft restoring a model
type RestoredHook interface {
	Restored(ctx context.Context, result *mongo.UpdateResult) error
}

// SavingHook is called before a model (new or existing) is saved to the database.
type SavingHook interface {
	Saving(ctx context.Context, cfg *FieldsConfig) error
}

// SavedHook is called after a model is saved to the database.
type SavedHook interface {
	Saved(ctx context.Context) error
}

// DeletingHook is called before a model is deleted
type DeletingHook interface {
	Deleting(ctx context.Context, cfg *FieldsConfig) error
}

// DeletedHook is called after a model is deleted
type DeletedHook interface {
	Deleted(ctx context.Context, result *mongo.DeleteResult) error
}
