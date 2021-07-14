/*
@Author: Weny Xu
@Date: 2021/06/03 18:11
*/

package odm

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/stretchr/testify/assert"
)

func TestOptions_NewCollection(t *testing.T) {
	opts := setupDefaultOpts()
	coll := opts.Coll(&Doc{})
	coll2 := opts.Coll(&Doc{})
	assert.Equal(t, coll, coll2)
}

func Test_prepareModels(t *testing.T) {
	validatedData := []*Doc{{Name: "test"}}
	invalidatedData := []Doc{{Name: "test"}}
	cfg := genModelFieldsInfo(&Doc{}).genModelConfig()
	output, err := prepareModels(context.Background(), cfg, validatedData)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	//assert.True(t, !validatedData[0].UpdatedAt.IsZero())
	output, err = prepareModels(context.Background(), cfg, invalidatedData)
	assert.NotNil(t, err)
	assert.Nil(t, output)
}

func TestCollection_Create(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	doc := Doc{
		Name: "weny",
		Age:  12,
	}
	id := primitive.NewObjectID()
	doc.SetID(id)
	err := opts.Coll(&doc).Create(context.TODO(), &doc)
	assert.Nil(t, err)
	assert.Equal(t, id, doc.ID)
	assert.False(t, doc.CreatedAt.IsZero())
	assert.False(t, doc.UpdatedAt.IsZero())
}

func TestCollection_CreateMany(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	docs := []*Doc{{Name: "weny"}, {Name: "leo"}}
	err := opts.Coll(&Doc{}).CreateMany(context.TODO(), &docs)
	assert.Nil(t, err)
	for _, doc := range docs {
		assert.False(t, doc.CreatedAt.IsZero())
		assert.False(t, doc.UpdatedAt.IsZero())
	}
}

func TestCollection_CreateMany_Fail(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	docs := []Doc{{Name: "weny"}, {Name: "leo"}}
	err := opts.Coll(&Doc{}).CreateMany(context.TODO(), &docs)
	assert.Equal(t, InvalidatedModels, err)
}

func TestCollection_UpdateOne(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	docs := seedDoc(opts)
	updatedAt := docs[0].UpdatedAt
	docs[0].Name = "weny updated"
	err := opts.Coll(&Doc{}).UpdateOne(context.TODO(), docs[0])
	assert.Nil(t, err)
	assert.NotEqual(t, docs[0].UpdatedAt, updatedAt)
	assert.Equal(t, "weny updated", docs[0].Name)
}

func TestCollection_Find(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	seedData := seedDoc(opts)
	var result []Doc
	err := opts.Coll(&Doc{}).Find(context.TODO(), bson.M{}, &result)
	assert.Nil(t, err)
	assert.Equal(t, len(seedData), len(result))
}

func TestCollection_FindOne(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	_ = seedDoc(opts)
	var result Doc
	err := opts.Coll(&Doc{}).FindOne(context.TODO(), bson.M{"name": "weny"}, &result)
	assert.Nil(t, err)
	assert.Equal(t, "weny", result.Name)
}

func TestCollection_SoftDeleteOne(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	data := seedDoc(opts)
	fmt.Println(data[0])
	err := opts.Coll(&Doc{}).SoftDeleteOne(context.TODO(), data[0])
	assert.Nil(t, err)
	fmt.Println(data[0])
	var result Doc
	err = opts.Coll(&Doc{}).FindOne(context.TODO(), bson.M{"_id": data[0].ID}, &result)
	assert.Equal(t, mongo.ErrNoDocuments, err)
}

func TestCollection_SoftDeleteMany(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	data := seedDoc(opts)
	fmt.Println(data[0])
	time.Sleep(100 * time.Millisecond)
	err := opts.Coll(&Doc{}).SoftDeleteMany(context.TODO(), bson.M{})
	assert.Nil(t, err)
	fmt.Println(data[0])
	var result []Doc
	err = opts.Coll(&Doc{}).Find(context.TODO(), bson.M{}, &result)
	assert.Equal(t, 0, len(result))
}

func TestCollection_RestoreOne(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	data := seedDoc(opts)
	// soft-deleting
	err := opts.Coll(&Doc{}).SoftDeleteOne(context.TODO(), data[0])
	assert.Nil(t, err)
	assert.False(t, data[0].DeletedAt.IsZero())
	// restoring
	err = opts.Coll(&Doc{}).RestoreOne(context.TODO(), data[0])
	// checking
	assert.Nil(t, err)
	assert.True(t, data[0].DeletedAt.IsZero())
}

func TestCollection_RestoreMany(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	data := seedDoc(opts)
	// soft-deleting
	err := opts.Coll(&Doc{}).SoftDeleteMany(context.TODO(), bson.M{})
	assert.Nil(t, err)
	var result []Doc
	// checking
	err = opts.Coll(&Doc{}).Find(context.TODO(), bson.M{}, &result)
	assert.Equal(t, 0, len(result))
	// restoring
	err = opts.Coll(&Doc{}).RestoreMany(context.TODO(), bson.M{})
	assert.Nil(t, err)
	// checking
	err = opts.Coll(&Doc{}).Find(context.TODO(), bson.M{}, &result)
	assert.Equal(t, len(data), len(result))
}

func TestCollection_DeleteOne(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	data := seedDoc(opts)

	// deleting
	err := opts.Coll(&Doc{}).DeleteOne(context.TODO(), data[0])
	assert.Nil(t, err)

}

func TestCollection_DeleteMany(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	_ = seedDoc(opts)

	// deleting
	err := opts.Coll(&Doc{}).DeleteMany(context.TODO(), bson.M{})
	assert.Nil(t, err)
}

func TestCollection_Count(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)
	data := seedDoc(opts)
	err := opts.Coll(&Doc{}).SoftDeleteOne(context.TODO(), data[0])
	assert.Nil(t, err)
	result, err := opts.Coll(&Doc{}).Count(context.TODO(), bson.M{})
	assert.Nil(t, err)
	assert.Equal(t, int64(len(data)-1), result)
}
