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

var (
	// ErrNotFound ...
	ErrNotFound = errors.New("not found")
)




type OperationType int64

const (
	//Insert ...
	Insert OperationType = iota
	//DeleteDoc ...
	DeleteDoc
	//Replace ...
	Replace
	//UpdateDoc ...
	UpdateDoc
	//Invalidate ...
	Invalidate
)

type ChangeEventCallbackFunc func(ce ChangeEventData)

type ChangeEventData struct {
	DocumentKey   primitive.ObjectID `bson:"documentKey" json:"documentKey"`
	OperationType string             `bson:"operationType" json:"operationType"`
	FullDocument  bson.Raw           `bson:"fullDocument" json:"fullDocument"`
}

func (ce *ChangeEventData) GetFullDocument(value interface{}) error {
	if ce.FullDocument != nil {
		err := bson.Unmarshal(ce.FullDocument, value)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (ce *ChangeEventData) GetID() string {
	if ce.DocumentKey.IsZero() {
		return ""
	}
	return ce.DocumentKey.Hex()
}

func (ce *ChangeEventData) GetType() OperationType {
	switch ce.OperationType {
	case "insert":
		return Insert
	case "delete":
		return DeleteDoc
	case "update":
		return UpdateDoc
	case "replace":
		return Replace
	default:
		return Invalidate
	}
}

func dbLayerWatch(ctx context.Context, collection *mongo.Collection, fullDocNeeded bool, callbacks ...ChangeEventCallbackFunc) error {
	//defer waitGroup.Done()
	pipeline := mongo.Pipeline{bson.D{{Key: "$addFields", Value: bson.D{{Key: "documentKey", Value: "$documentKey._id"}}}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "documentKey", Value: 1}, {Key: "operationType", Value: 1}, {Key: "fullDocument", Value: 1}}}}}

	var streamOptions *options.ChangeStreamOptions
	if fullDocNeeded {
		streamOptions = options.ChangeStream().SetFullDocument(options.UpdateLookup)
	} else {
		streamOptions = options.ChangeStream().SetFullDocument(options.Default)
	}
	stream, err := collection.Watch(ctx, pipeline, streamOptions)
	if err != nil {
		fmt.Printf("Error: encountered error watching for change stream data, error is %v\n", err)
		return err
	}
	defer func(stream *mongo.ChangeStream, ctx context.Context) {
		err := stream.Close(ctx)
		if err != nil {
			fmt.Println("Client.Watch>Stream.Close err=>", err)
		}
	}(stream, ctx)
	for stream.Next(ctx) {
		var data ChangeEventData
		if err := stream.Decode(&data); err != nil {
			fmt.Printf("Error: failed to extract change event data from change stream, error is %v\n", err)
			continue
		}

		for _, callback := range callbacks {
			go func(f func(ev ChangeEventData), e ChangeEventData) {
				f(e)
			}(callback, data)
		}
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return nil
}

type ChangeEventModelCallbackFunc func(id string, optType OperationType, m *Model)

func modelLayerWatch(ctx context.Context, client *mongo.Database, fullDocNeeded bool, callback ChangeEventModelCallbackFunc) error {
	collection := client.Collection("episodes")

	err := dbLayerWatch(ctx, collection, fullDocNeeded, func(callback ChangeEventModelCallbackFunc) ChangeEventCallbackFunc {
		return func(ce ChangeEventData) {
			var val *Model
			if fullDocNeeded {
				var temp Model
				err := ce.GetFullDocument(&temp)
				if err == nil {
					val = &temp
				}
			}

			callback(ce.DocumentKey.Hex(), ce.GetType(), val)
		}
	}(callback))
	return err
}

func callback() ChangeEventModelCallbackFunc {
	return func(id string, optType OperationType, m *Model) {
		if m == nil {
			fmt.Println("empty doc on client: ", id, optType)
		} else {
			if optType == DeleteDoc {
				fmt.Println(optType)
			}
			fmt.Println("Final method - id: ", id, "name: ", m.Name, "age: ", m.Age, "modelId: ", m.ID, "optType: ", optType)
		}
	}
}