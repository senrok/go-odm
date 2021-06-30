/*
@Author: Weny Xu
@Date: 2021/06/03 17:29
*/

package odm

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/joho/godotenv"
)

var (
	MONGODB_URL     string
	MONGODB_DB_NAME = "odm_test"
	CI_ENV          = false
)

func loadConfig() {
	if os.Getenv("CI") != "true" {
		if err := godotenv.Load(".env_local"); err != nil {
			println("loading .env_local failed: ", err.Error())
		}
	} else {
		CI_ENV = true
	}
	MONGODB_URL = os.Getenv("MONGODB_URL")
}

func setupDefaultOpts() *Options {
	opts, err := DefaultOpts(SetDatabase(MONGODB_URL, MONGODB_DB_NAME))
	if err != nil {
		panic(err)
	}
	return opts
}

func resetDB(opts *Options) {
	_, err := opts.Coll(&Doc{}).Collection.DeleteMany(context.TODO(), bson.M{})
	PanicErr(err)
}

func seedDoc(opts *Options) []*Doc {
	docs := []*Doc{{Name: "weny"}, {Name: "leo"}}
	err := opts.Coll(&Doc{}).CreateMany(context.TODO(), &docs)
	PanicErr(err)
	return docs
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Doc struct {
	DefaultModel `bson:",inline"`

	Name string `bson:"name"`
	Age  int    `bson:"age"`
}
