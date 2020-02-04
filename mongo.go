package main

import (
  "log"
  "fmt"
  "time"
  "context"
  "github.com/crackhd/env"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoBackend struct {
  env    *env.Env
  client *mongo.Client
  dbName string
}

func NewMongoBackend(env *env.Env) (*MongoBackend, error) {

  mongoHost, _ := env.Get("DB_HOST", "localhost")
  mongoPort, _ := env.Get("DB_PORT", "27017")
  dbName, _ := env.Get("DB_NAME", "changeme-dev")

  // Set client options
  mongos := fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)
  clientOptions := options.Client().ApplyURI(mongos)

  // Connect to MongoDB

  log.Println("Connecting to MongoDB at", mongoHost, mongoPort)
  client, err := mongo.Connect(context.TODO(), clientOptions)

  if err != nil {
    return nil, err
  }

  // Check the connection
  err = client.Ping(context.TODO(), nil)

  for i := 0; err != nil && i < 10; i++ {
    fmt.Println("MongoDB connection fail:", mongoHost, mongoPort, "-", err.Error(), ", retrying...")
    time.Sleep(5 * time.Second)
    err = client.Ping(context.TODO(), nil)
  }

  if err != nil {
    log.Println("MongoDB connection failed: ", err.Error())
  }

  log.Println("Connected to MongoDB.")

  return &MongoBackend{env, client, dbName}, nil
}

func (b *MongoBackend) Collection(name string) *mongo.Collection {
  return b.client.Database(b.dbName).Collection(name)
}
