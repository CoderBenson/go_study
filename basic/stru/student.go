package stru

import "fmt"

type Student struct {
	Age  uint8
	Name string
	Man  bool
	Mark float32
}

func NewStudent(age uint8, name string, man bool, mark float32) *Student {
	return &Student{
		Age:  age,
		Name: name,
		Man:  man,
		Mark: mark,
	}
}

func (stu *Student) Eat(food string) {
	fmt.Printf("%s eat %s", stu.Name, food)
}

type OtherSutdent struct {
	Student
	SecondName string
}

func NewOtherStudent(student Student, secondName string) *OtherSutdent {
	return &OtherSutdent{
		Student:    student,
		SecondName: secondName,
	}
}

func NewOtherStudentDetail(
	age uint8, name string, man bool, mark float32,
	secondName string,
) *OtherSutdent {
	return NewOtherStudent(*NewStudent(age, name, man, mark), secondName)
}
