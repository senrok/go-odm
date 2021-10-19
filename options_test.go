/*
@Author: Weny Xu
@Date: 2021/06/03 17:28
*/

package odm

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func ExampleDefaultOpts() {
	_, err := DefaultOpts(SetDatabase(MONGODB_URL, MONGODB_DB_NAME))
	if err != nil {
		panic(err)
	}
}

func ExampleNewOpts() {
	_, err := NewOpts(func(opts *Options) error {
		url := "your connection string"
		c, err := mongo.NewClient(options.Client().ApplyURI(url))
		if err != nil {
			return err
		}
		// specify the database for opts
		opts.db = c.Database("your database name")
		return nil
	})
	if err != nil {
		panic(err)
	}
}
