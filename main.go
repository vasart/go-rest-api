package main

import (
	"log"

	"github.com/vasart/go-rest-api/db"
)

func main() {
	mgoSession, err := db.NewSession("127.0.0.1:27017")
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}
	defer mgoSession.Close()

	u := db.NewUserService(mgoSession.Copy(), "go-rest-api", "user")
	s := NewServer(u)

	s.Start()
}
