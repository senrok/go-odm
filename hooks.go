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
	SoftDeleting(ctx context.Context, cfg *FieldsConfig) error
}

// SoftDeletedHook is called after soft deleting a model
type SoftDeletedHook interface {
	SoftDeleted(ctx context.Context, result *mongo.UpdateResult) error
}

// RestoringHook is called before soft restoring a model
type RestoringHook interface {
	Restoring(ctx context.Context, cfg *FieldsConfig) error
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

type hookRunner func(ctx context.Context, cfg *FieldsConfig, model IModel) error

func restoringHook(ctx context.Context, cfg *FieldsConfig, model IModel) error {
	if hook, ok := model.(RestoringHook); ok {
		if err := hook.Restoring(ctx, cfg); err != nil {
			return err
		}
	}
	return nil
}

func restoredHook(updateResult *mongo.UpdateResult) hookRunner {
	return func(ctx context.Context, cfg *FieldsConfig, model IModel) error {
		if hook, ok := model.(RestoredHook); ok {
			if err := hook.Restored(ctx, updateResult); err != nil {
				return err
			}
		}
		return nil
	}
}

func softDeletingHook(ctx context.Context, cfg *FieldsConfig, model IModel) error {
	if hook, ok := model.(SoftDeletingHook); ok {
		if err := hook.SoftDeleting(ctx, cfg); err != nil {
			return err
		}
	}
	return nil
}

func softDeletedHook(updateResult *mongo.UpdateResult) hookRunner {
	return func(ctx context.Context, cfg *FieldsConfig, model IModel) error {
		if hook, ok := model.(SoftDeletedHook); ok {
			if err := hook.SoftDeleted(ctx, updateResult); err != nil {
				return err
			}
		}
		return nil
	}
}

func updatingHook(ctx context.Context, cfg *FieldsConfig, model IModel) error {
	if hook, ok := model.(UpdatingHook); ok {
		if err := hook.Updating(ctx); err != nil {
			return err
		}
	}
	return nil
}

func updatedHook(updateResult *mongo.UpdateResult) hookRunner {
	return func(ctx context.Context, cfg *FieldsConfig, model IModel) error {
		if hook, ok := model.(SavedHook); ok {
			if err := hook.Saved(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

func savedHook(ctx context.Context, cfg *FieldsConfig, model IModel) error {
	if hook, ok := model.(SavedHook); ok {
		if err := hook.Saved(ctx); err != nil {
			return err
		}
	}
	return nil
}

func savingHook(ctx context.Context, cfg *FieldsConfig, model IModel) error {
	return model.Saving(ctx, cfg)
}

func creatingHook(ctx context.Context, cfg *FieldsConfig, model IModel) error {
	return model.Creating(ctx, cfg)
}

func createdHook(ctx context.Context, cfg *FieldsConfig, model IModel) error {
	if hook, ok := model.(CreatedHook); ok {
		if err := hook.Created(ctx); err != nil {
			return err
		}
	}
	return nil
}

func deletingHook(ctx context.Context, cfg *FieldsConfig, model IModel) error {
	if hook, ok := model.(DeletingHook); ok {
		if err := hook.Deleting(ctx, cfg); err != nil {
			return err
		}
	}
	return nil
}

func deletedHook(result *mongo.DeleteResult) hookRunner {
	return func(ctx context.Context, cfg *FieldsConfig, model IModel) error {
		if hook, ok := model.(DeletedHook); ok {
			if err := hook.Deleted(ctx, result); err != nil {
				return err
			}
		}
		return nil
	}
}

func modelHooksRunnerExecutor(ctx context.Context, cfg *FieldsConfig, model IModel, runners ...hookRunner) error {
	for _, runner := range runners {
		err := runner(ctx, cfg, model)
		if err != nil {
			return err
		}
	}
	return nil
}

func modelsHooksRunnerExecutor(ctx context.Context, cfg *FieldsConfig, models IModels, runners ...hookRunner) error {
	for _, model := range models {
		err := modelHooksRunnerExecutor(ctx, cfg, model, runners...)
		if err != nil {
			return err
		}
	}
	return nil
}
