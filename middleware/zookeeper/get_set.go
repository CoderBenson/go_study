package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/CoderBenson/go_study/tool/str"
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

func ObtainConn() *zk.Conn {
	return conn_pool.Get().(*zk.Conn)
}

func RecycleConn(conn *zk.Conn) {
	conn_pool.Put(conn)
}

func DoWithConn(fun func(*zk.Conn)) {
	conn := ObtainConn()
	defer RecycleConn(conn)
	fun(conn)
}

func create(path, value string) (string, error) {
	var data = []byte(value)
	var flags int32 = 0
	acls := zk.WorldACL(zk.PermAll)
	conn := ObtainConn()
	defer RecycleConn(conn)

	res, err := conn.Create(path, data, flags, acls)
	if err != nil {
		return "", err
	}
	return res, nil
}

func get(path string) (string, error) {
	conn := ObtainConn()
	result, _, err := conn.Get(path)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func exists(path string) (bool, error) {
	conn := ObtainConn()
	defer RecycleConn(conn)
	r, _, err := conn.Exists(path)
	return r, err
}

func set(path, value string) error {
	conn := ObtainConn()
	defer RecycleConn(conn)
	var flags int32 = 0
	_, err := conn.Set(path, []byte(value), flags)
	return err
}

const (
	Operate_Create = "create"
	Operate_Get    = "get"
	Operate_Set    = "set"
	Operate_Exist  = "exist"
	Operate_Child  = "child"
)

func child(path string) ([]string, error) {
	conn := ObtainConn()
	defer RecycleConn(conn)
	cs, _, err := conn.Children(path)
	return cs, err
}

func main() {
	operateStr := &os.Args[1]
	pathStr := &os.Args[2]
	var valueStr *string
	if *operateStr == Operate_Create || *operateStr == Operate_Set {
		valueStr = &os.Args[3]
	}

	if str.EmptyTrim(*operateStr) {
		panic(errors.New("operate should be input"))
	}
	switch *operateStr {
	case Operate_Create:
		res, err := create(*pathStr, *valueStr)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(res)
			fmt.Printf("create node successful:%s->%s\n", *pathStr, *valueStr)
		}
	case Operate_Get:
		if result, err := get(*pathStr); err != nil {
			panic(err)
		} else {
			fmt.Printf("get from zookeeper successful:%s\n", result)
		}
	case Operate_Exist:
		if r, err := exists(*pathStr); err != nil {
			panic(err)
		} else {
			fmt.Printf("%s %v exists in zookeeper\n", *pathStr, r)
		}
	case Operate_Set:
		if err := set(*pathStr, *valueStr); err != nil {
			panic(err)
		} else {
			fmt.Printf("set node successful:%s->%s\n", *pathStr, *valueStr)
		}
	case Operate_Child:
		if cs, err := child(*pathStr); err != nil {
			panic(err)
		} else {
			fmt.Printf("get %s's child from zookeeper successful:%v\n", *pathStr, cs)
		}
	default:
		panic(fmt.Errorf("%s operate not support", *operateStr))
	}
}
