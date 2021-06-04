/*
@Author: Weny Xu
@Date: 2021/06/03 17:28
*/

package odm

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	loadConfig()
	os.Exit(m.Run())
}

func TestSetDatabase(t *testing.T) {
	opts, err := DefaultOpts(SetDatabase(MONGODB_URL, MONGODB_DB_NAME))
	if err != nil {
		panic(err)
	}
	assert.NotNil(t, opts.db)
}
