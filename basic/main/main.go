package main

import (
	"fmt"

	"github.com/CoderBenson/go_study/basic/inface"
	"github.com/CoderBenson/go_study/basic/stru"
)

func simpleStudent() {
	stu := stru.Student{}
	stu.Age = 10
	stu1 := stu
	stu1.Age = 20
	fmt.Printf("stu:%v,stu1:%v\n", stu, stu1)
	pStu1 := &stu
	pStu1.Age = 20
	fmt.Printf("stu:%v, pStu1:%v\n", stu, *pStu1)

	pStu2 := &stu1
	fmt.Println(stu == stu1)
	fmt.Println(stu == *pStu1)
	fmt.Println(pStu1 == pStu2)
}

func extend() {
	stu := stru.Student{Name: "xiaoming"}
	other := stru.OtherSutdent{Student: stu, SecondName: "Jack"}
	fmt.Printf("%+v\n", other)
	fmt.Printf("other.Name=%s\n", other.Name)
}

func usePoint() {
	stu := *stru.NewStudent(20, "xiaoming", true, 20)
	pStu1 := &stu
	pStu2 := pStu1
	fmt.Println(pStu1 == pStu2)
	pStu1 = stru.NewStudent(20, "xiaohong", false, 30)
	fmt.Println(pStu1 == pStu2)
}

func main() {
	inface.TestInterface()
}
