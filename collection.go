/*
@Author: Weny Xu
@Date: 2021/06/02 22:42
*/

package odm

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
