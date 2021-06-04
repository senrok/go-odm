/*
@Author: Weny Xu
@Date: 2021/06/03 18:11
*/

package odm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptions_NewCollection(t *testing.T) {
	opts := setupDefaultOpts()
	coll := opts.Coll(&Doc{})
	coll2 := opts.Coll(&Doc{})
	assert.Equal(t, coll, coll2)
}
