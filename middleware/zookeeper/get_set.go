package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
)

var conn_pool = sync.Pool{New: func() any {
	hosts := []string{"127.0.0.1:2181"}
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		panic(err)
	}
	return conn
}}

func create(path string) {
	var data = []byte("test value")
	var flags int32 = 0
	acls := zk.WorldACL(zk.PermAll)
	conn := conn_pool.Get().(*zk.Conn)

	s, err := conn.Create(path, data, flags, acls)
	if err != nil {
		panic(err)
	}
	fmt.Printf("create succssful:%s\n", s)
}

func get(path string) {
	conn := conn_pool.Get().(*zk.Conn)
	result, _, err := conn.Get(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("get successful:%s\n", string(result))
}

func main() {
	path := "/test"
	// create(path)
	get(path)
}
