/*
@Author: Weny Xu
@Date: 2021/06/02 22:57
*/

package odm

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

type TagTest struct {
	DefaultModel `bson:",inline"`
}

func TestGenModelConfig(t *testing.T) {
	data := TagTest{}
	cfg := genModelFieldsInfo(data).genModelConfig()
	_ = data.Creating(context.TODO(), cfg)
	assert.False(t, time.Time{}.Equal(data.CreatedAt))
	_ = data.Saving(context.TODO(), cfg)
	assert.False(t, time.Time{}.Equal(data.UpdatedAt))
	_ = data.Deleting(context.TODO(), cfg)
	assert.False(t, time.Time{}.Equal(data.DeletedAt))
}

func TestModelTag(t *testing.T) {
	data := TagTest{}
	rv := reflect.TypeOf(data)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	//fmt.Println(rv.FieldByIndex([]int{0, 1, 0}))
	var i []int
	var fi FieldsInfo
	getModelFieldsInfo(rv, &fi, i)
	fmt.Println(fi)
}
