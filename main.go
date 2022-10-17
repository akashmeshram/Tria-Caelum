package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/allegro/bigcache/v3"
	cacheStoreCache "github.com/eko/gocache/cache"
	cacheStore "github.com/eko/gocache/store"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


//Model ChangeEventCallbackFunc func(id string, m Model) error
type Model struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name,omitempty"`
	Age  string             `bson:"age,omitempty"`
}



func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("url"))
	if err != nil {
		panic(err)
	} else {
		fmt.Println("mongo connected")
	}
	c, err := BigCache(10 * time.Minute)
	if err != nil {
		fmt.Println("Err in cache")
	}
	id, _ := primitive.ObjectIDFromHex("62e74de5413043a5eea730f0")
	temp := Model{
		ID:   id,
		Name: "hhdj",
		Age:  "78",
	}
	err = c.Set(context.Background(), "62e75de5413043a5eea730f0", temp)
	if err != nil {
		fmt.Println("set error: ", err)
	}
	fmt.Println("Data saved to cache")

	var res Model
	err = c.Get(context.Background(), "62e75de5413043a5eea730f0", &res)

	if err != nil {
		fmt.Println("error 199:", err)
	} else {
		fmt.Println("hey:", res.ID, res.Name, res.Age)
	}

	defer client.Disconnect(context.TODO())

	
}
