package main

import (
	"kvstore/server"
	"kvstore/store"
)

func main() {
	db := store.InitInMemDb()
	proto := &server.RedisProtocol{}
	store.Init(db)
	srv := server.CreateServer(8080, proto)
	srv.Start()
}
