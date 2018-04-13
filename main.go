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

	// Insert
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

	// Create Index
	indexName, err := c.Indexes().CreateOne(
		context.Background(),
		mgo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.Int32("time", 1),
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

	// Find
	cur, err := c.Find(context.Background(), bson.NewDocument(
		bson.EC.SubDocument("hello", bson.NewDocument(
			bson.EC.Boolean("$exists", true),
		)),
	))
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.Background()) {
		elem := bson.NewDocument()

		if err := cur.Decode(elem); err != nil {
			log.Fatal(err)
		}

		fmt.Println(elem)
	}
	cur.Close(context.Background())

	// Aggregate
	cur, err = c.Aggregate(context.Background(), []*bson.Document{
		bson.NewDocument(bson.EC.SubDocument("$match", bson.NewDocument(
			bson.EC.String("name", "Parham"),
		))),
		bson.NewDocument(bson.EC.SubDocument("$sort", bson.NewDocument(
			bson.EC.Int32("time", 1),
		))),
	})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.Background()) {
		elem := bson.NewDocument()

		if err := cur.Decode(elem); err != nil {
			log.Fatal(err)
		}

		fmt.Println(elem)
	}
	cur.Close(context.Background())

}
