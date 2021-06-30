/*
@Author: Weny Xu
@Date: 2021/06/03 10:48
*/

package odm

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collections sync.Map

type Options struct {
	db          *mongo.Database
	timeout     time.Duration
	collections *sync.Map
}

type Option func(opts *Options) error

//Options Clone Current Options
func (o *Options) Clone() *Options {
	// clone collections
	collectionsCopy := &sync.Map{}
	o.collections.Range(func(key, value interface{}) bool {
		collectionsCopy.Store(key, value)
		return true
	})
	return &Options{
		db:          o.db,
		timeout:     o.timeout,
		collections: collectionsCopy,
	}
}

// SetDatabase setup database by connection url and db name
func SetDatabase(url string, dbName string, opts ...*options.ClientOptions) Option {
	return func(opt *Options) error {
		c, err := mongo.NewClient(
			append(
				[]*options.ClientOptions{
					options.Client().ApplyURI(url),
				}, opts...,
			)...,
		)
		if err != nil {
			return err
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = c.Connect(ctx)
		if err != nil {
			return err
		}

		opt.db = c.Database(dbName)
		return nil
	}
}

var defaultOptions = &Options{
	db:          nil,
	timeout:     time.Second * 10,
	collections: &sync.Map{},
}

// DefaultOpts returns a new Options
func NewOpts(opts ...Option) (*Options, error) {
	opt := &Options{
		db:          nil,
		timeout:     time.Second * 10,
		collections: &sync.Map{},
	}
	for _, o := range opts {
		if o != nil {
			_ = o(opt)
		}
	}
	return opt, nil
}

// DefaultOpts returns a Options cloned from defaultOptions
func DefaultOpts(opts ...Option) (*Options, error) {
	clone := defaultOptions.Clone()
	for _, o := range opts {
		if o != nil {
			_ = o(clone)
		}
	}
	return clone, nil
}
