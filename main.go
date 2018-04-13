/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 13-04-2018
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
)

func main() {
	fmt.Println("vim-go")

	client, err := mgo.NewClient("mongodb://127.0.0.1:27017")
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("test")
	c := db.Collection("nice")

	s := struct {
		Name   string
		Family string
		Time   time.Time
	}{
		Name:   "Parham",
		Family: "Alvani",
		Time:   time.Now(),
	}
	fmt.Println(s)

	res, err := c.InsertOne(context.Background(), s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.InsertedID)

	indexName, err := c.Indexes().CreateOne(
		context.Background(),
		mgo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.Int32("foo", 1),
			),
			Options: bson.NewDocument(
				bson.EC.Boolean("unique", true),
			),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(indexName)

	cur, err := c.Find(context.Background(), bson.NewDocument(
		bson.EC.SubDocument("hello", bson.NewDocument(
			bson.EC.Boolean("$exists", true),
		)),
	))
	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		elem := bson.NewDocument()

		if err := cur.Decode(elem); err != nil {
			log.Fatal(err)
		}

		fmt.Println(elem)
	}
}
