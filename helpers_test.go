/*
@Author: Weny Xu
@Date: 2021/06/03 17:29
*/

package odm

import (
	"github.com/joho/godotenv"
	"os"
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

type Doc struct {
	DefaultModel `bson:",inline"`

	Name string `bson:"name"`
	Age  int    `bson:"age"`
}
