/*
@Author: Weny Xu
@Date: 2021/06/03 18:11
*/

package odm

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
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
