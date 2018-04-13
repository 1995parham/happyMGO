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

	mgo "github.com/mongodb/mongo-go-driver/mongo"
)

func main() {
	fmt.Println("vim-go")

	client, err := mgo.NewClient("mongodb://127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("test")
	c := db.Collection("hello")

	res, err := c.InsertOne(context.Background(), map[string]string{"hello": "world"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.InsertedID)
}
