package main

import (
	"context"
	"time"

	"github.com/sunshineplan/utils/database/mongodb"
)

var mongo mongodb.Config

func query(metadata string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Open()
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database(mongo.Database).Collection(mongo.Collection)

	return collection.FindOne(ctx, map[string]interface{}{"_id": metadata}).Decode(data)
}
