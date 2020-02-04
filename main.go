package main

import (
  "log"
)

var (
  env Env
  mgb *MongoBackend
)

func init() {
  var err error

	if err = env.LoadFile(); err != nil {
		log.Println("WARN: Error reading .env file:", err.Error())
	}

  mgb, err = NewMongoBackend(&env)
  if err != nil {
    panic(err)
  }
}

func main() {
  log.Println("changeme works!")
}
