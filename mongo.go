package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mc mongoConfig

type mongoConfig struct {
	Server     string
	Port       int
	Database   string
	Collection string
	Username   string
	Password   string
}

func query(metadata string) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", mc.Username, mc.Password, mc.Server, mc.Port, mc.Database)))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)
	collection := client.Database(mc.Database).Collection(mc.Collection)

	var result bson.M
	if err = collection.FindOne(ctx, bson.M{"_id": metadata}).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func dump() string {
	tmpfile, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	tmpfile.Close()

	args := []string{}
	args = append(args, fmt.Sprintf("-h%s:%d", mc.Server, mc.Port))
	args = append(args, fmt.Sprintf("-d%s", mc.Database))
	args = append(args, fmt.Sprintf("-c%s", mc.Collection))
	args = append(args, fmt.Sprintf("-u%s", mc.Username))
	args = append(args, fmt.Sprintf("-p%s", mc.Password))
	args = append(args, "--gzip")
	args = append(args, fmt.Sprintf("--archive=%s", tmpfile.Name()))
	cmd := exec.Command("mongodump", args...)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return tmpfile.Name()
}
