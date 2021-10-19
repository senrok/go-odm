/*
@Author: Weny Xu
@Date: 2021/10/19 20:28
*/

package main

import (
	"context"
	"fmt"
	"github.com/senrok/go-odm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type Doc struct {
	odm.DefaultModel `bson:",inline"`

	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func main() {
	opts, err := odm.DefaultOpts(odm.SetDatabase(os.Getenv("URL"), os.Getenv("DB_NAME")))
	if err != nil {
		panic(err)
	}
	coll := opts.Coll(&Doc{})
	// reset
	coll.Collection.DeleteOne(context.Background(), bson.M{})
	fmt.Println("inserting a document")
	weny := &Doc{
		Name: "weny",
		Age:  12,
	}
	// Create
	err = coll.Create(context.Background(), weny)
	if err != nil {
		fmt.Println("failed to insert the document")
		return
	}
	// Create many
	docs := []*Doc{{Name: "weny"}, {Name: "leo"}}
	fmt.Println("inserting multi documents")
	err = coll.CreateMany(context.Background(), &docs)
	if err != nil {
		fmt.Println("failed to insert multi documents")
		return
	}

	// Update
	weny.Name = "updated weny"
	fmt.Println("updating a document")
	err = coll.UpdateOne(context.Background(), weny)
	if err != nil {
		fmt.Println("failed to update the documents")
		return
	}

	fmt.Println("finding a document")
	var result []Doc
	err = coll.Find(context.Background(), bson.M{}, &result)
	if err != nil {
		fmt.Println("failed to find documents")
		return
	}
	fmt.Printf("found: %v\n", result)

	fmt.Println("finding leo")
	var leo Doc
	// Find One
	err = coll.FindOne(context.Background(), bson.M{"name": "leo"}, &leo)
	if err != nil {
		fmt.Println("failed to find the first documents")
	}
	fmt.Printf("found: %v\n", leo)

	fmt.Println("soft-delete leo")
	// Soft Delete One
	err = coll.SoftDeleteOne(context.Background(), &leo)
	if err != nil {
		fmt.Println("failed to soft-delete document")
	}
	fmt.Println("trying find the leo")
	var afterSoftDeletedLeo Doc

	err = coll.FindOne(context.Background(), bson.M{"name": "leo"}, &afterSoftDeletedLeo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("soft deleted!")
		} else {
			fmt.Println("failed to find the first documents")
		}
	}
	fmt.Printf("found: %v\n", leo)

	// Soft Delete Many
	err = coll.SoftDeleteMany(context.Background(), bson.M{})
	if err != nil {
		fmt.Println("failed to soft delete multi documents")
	}
	fmt.Println("restoring leo")
	// Restore One
	err = coll.RestoreOne(context.Background(), &leo)
	if err != nil {
		fmt.Println("failed to restore a documents")
	}

	fmt.Println("restoring rest data")
	// Restore Many
	err = coll.RestoreMany(context.Background(), bson.M{})
	if err != nil {
		fmt.Println("failed to restore multi documents")
	}

	fmt.Println("terminated.")

}
