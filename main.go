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

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// create mongodb connection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// connect to the mongodb
	ctxc, donec := context.WithTimeout(context.Background(), 10*time.Second)
	defer donec()
	if err := client.Connect(ctxc); err != nil {
		logrus.Fatalf("db connection error: %s", err)
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
		mongo.IndexModel{
			Keys: bson.M{
				"time": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(indexName)

	// Find
	cur, err := c.Find(context.Background(), bson.M{
		"hello": bson.M{
			"$exists": true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.Background()) {
		var elem bson.D

		if err := cur.Decode(&elem); err != nil {
			log.Fatal(err)
		}

		fmt.Println(elem)
	}
	cur.Close(context.Background())

	// Aggregate
	cur, err = c.Aggregate(context.Background(), bson.A{
		bson.M{
			"$match": bson.M{
				"name": "Parham",
			},
		},
		bson.M{
			"$sort": bson.M{
				"time": 1,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.Background()) {
		var elem bson.D

		if err := cur.Decode(&elem); err != nil {
			log.Fatal(err)
		}

		fmt.Println(elem)
	}
	cur.Close(context.Background())
}
