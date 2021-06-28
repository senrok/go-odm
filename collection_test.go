/*
@Author: Weny Xu
@Date: 2021/06/03 18:11
*/

package odm

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
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

func MockAssertionFn(t interface{}) {
	rv := reflect.ValueOf(t)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		fmt.Println(rv.Kind(), rv.Type(), rv)
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			if hook, ok := rv.Index(i).Interface().(IModel); ok {
				fmt.Println("run hooks")
				hook.Saving(context.Background(), nil)
			}
		}
	default:
		fmt.Print("invalidated input")
	}
	fmt.Println(rv)
}

func TestReflectAssert(t *testing.T) {
	data := []*Doc{{Name: "test"}}

	MockAssertionFn(data)

}
