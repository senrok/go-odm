package odm

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

// Note: to run Transaction tests, the MongoDB daemon must run as replica set, not as a standalone daemon.
// To convert it [see this](https://docs.mongodb.com/manual/tutorial/convert-standalone-to-replica-set/)
func TestTransactionCommit(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)

	d := &Doc{Name: "check", Age: 10}

	err := opts.TransactionWithCtx(context.TODO(), func(session mongo.Session, sc mongo.SessionContext) error {

		err := opts.Coll(d).Create(sc, d)

		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})

	assert.Nil(t, err)
	count, err := opts.Coll(d).CountDocuments(context.TODO(), bson.M{})

	assert.Nil(t, err)
	require.Equal(t, int64(1), count)
}

func TestTransactionAbort(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)

	//seed()

	d := &Doc{Name: "check", Age: 10}

	err := opts.TransactionWithCtx(context.TODO(), func(session mongo.Session, sc mongo.SessionContext) error {

		err := opts.Coll(d).Create(sc, d)

		if err != nil {
			return err
		}

		return session.AbortTransaction(sc)
	})

	assert.Nil(t, err)
	count, err := opts.Coll(d).CountDocuments(context.TODO(), bson.M{})

	assert.Nil(t, err)
	require.Equal(t, int64(0), count)
}

func TestTransactionWithCtx(t *testing.T) {
	opts := setupDefaultOpts()
	resetDB(opts)

	//seed()

	d := &Doc{Name: "check", Age: 10}

	err := opts.TransactionWithCtx(context.TODO(), func(session mongo.Session, sc mongo.SessionContext) error {

		err := opts.Coll(d).Create(sc, d)

		if err != nil {
			return err
		}

		return session.AbortTransaction(sc)
	})

	assert.Nil(t, err)
	count, err := opts.Coll(d).CountDocuments(context.TODO(), bson.M{})

	assert.Nil(t, err)
	require.Equal(t, int64(0), count)
}
