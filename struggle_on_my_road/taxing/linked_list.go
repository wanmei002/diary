package main

import (
	"fmt"
	"math/rand"
)

type Student struct {
	Name  string
	Age   int
	Score float32
	next  *Student
}

func insertHead(p **Student){
	for i:=0; i<10; i++ {
		stu := Student{
			Name: fmt.Sprintf("stu%d", i),
			Age:  rand.Intn(100),
			Score: rand.Float32() * 100,
		}
		stu.next = *p
		*p = &stu
	}
}

func trans(p *Student){
	for p != nil {
		fmt.Println(*p)
		p = p.next
	}
}

func main(){

	var head *Student = new(Student)

	head.Name = "zyn"
	head.Age = 3
	head.Score = 100.0

	insertHead(&head)

	trans(head)
}