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
	RecycleConn(conn)
	result, _, err := conn.Get(path)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func getW(path string) (string, <-chan zk.Event, error) {
	conn := ObtainConn()
	defer RecycleConn(conn)
	r, _, rChan, err := conn.GetW(path)
	return string(r), rChan, err
}

func exists(path string) (bool, error) {
	conn := ObtainConn()
	defer RecycleConn(conn)
	r, _, err := conn.Exists(path)
	return r, err
}

func existsw(path string) (bool, <-chan zk.Event, error) {
	conn := ObtainConn()
	defer RecycleConn(conn)
	r, _, rChan, err := conn.ExistsW(path)
	return r, rChan, err
}

func set(path, value string) (string, error) {
	conn := ObtainConn()
	defer RecycleConn(conn)
	old, state, err := conn.Get(path)
	if err != nil {
		return "", nil
	}
	_, err = conn.Set(path, []byte(value), state.Version)
	return string(old), err
}

func child(path string) ([]string, error) {
	conn := ObtainConn()
	defer RecycleConn(conn)
	cs, _, err := conn.Children(path)
	return cs, err
}

func childw(path string) ([]string, <-chan zk.Event, error) {
	conn := ObtainConn()
	defer RecycleConn(conn)

	cs, _, csChann, err := conn.ChildrenW(path)
	return cs, csChann, err
}

const (
	Operate_Create = "create"
	Operate_Get    = "get"
	Operate_Getw   = "getw"
	Operate_Set    = "set"
	Operate_Exist  = "exist"
	Operate_Existw = "existw"
	Operate_Child  = "child"
	Operate_Childw = "childw"
)

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
	case Operate_Getw:
		if r, rChan, err := getW(*pathStr); err != nil {
			panic(err)
		} else {
			fmt.Printf("get from zookeeper successful:%s\n", r)
			event := <-rChan
			if err = event.Err; err != nil {
				panic(event.Err)
			} else if event.Type == zk.EventNodeDataChanged {
				r, err = get(*pathStr)
				if err != nil {
					panic(err)
				} else {
					fmt.Printf("getw channel watch:%s->%s\n", *pathStr, r)
				}
			}
		}
	case Operate_Exist:
		if r, err := exists(*pathStr); err != nil {
			panic(err)
		} else {
			fmt.Printf("%s %v exists in zookeeper\n", *pathStr, r)
		}
	case Operate_Existw:
		if r, rChan, err := existsw(*pathStr); err != nil {
			panic(err)
		} else {
			fmt.Printf("%s %v existsw in zookeeper\n", *pathStr, r)
			event := <-rChan
			if err = event.Err; err != nil {
				panic(err)
			} else if event.Type == zk.EventNodeCreated {
				str, err := get(*pathStr)
				if err != nil {
					panic(err)
				} else {
					fmt.Printf("existsw in zookeeper channel watch:%s->%s\n", *pathStr, str)
				}
			}
		}
	case Operate_Set:
		if old, err := set(*pathStr, *valueStr); err != nil {
			panic(err)
		} else {
			fmt.Printf("set node(%s) successful:%s->%s\n", old, *pathStr, *valueStr)
		}
	case Operate_Child:
		if cs, err := child(*pathStr); err != nil {
			panic(err)
		} else {
			fmt.Printf("get %s's child from zookeeper successful:%v\n", *pathStr, cs)
		}
	case Operate_Childw:
		if cs, csChan, err := childw(*pathStr); err != nil {
			panic(err)
		} else {
			fmt.Printf("get %s's child from zookeeper successful:%v\n", *pathStr, cs)
			event := <-csChan
			if err = event.Err; err != nil {
				panic(err)
			} else if event.Type == zk.EventNodeChildrenChanged {
				cs, err = child(*pathStr)
				if err != nil {
					panic(err)
				} else {
					fmt.Printf("get %s child from zookeeper channel watch:%v\n", *pathStr, cs)
				}
			}
		}
	default:
		panic(fmt.Errorf("%s operate not support", *operateStr))
	}
}
