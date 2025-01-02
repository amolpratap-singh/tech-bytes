package main

import (
	"fmt"
	"strings"

	array "learnGoWithTest/internal/array"
	hello "learnGoWithTest/internal/hello_world"
	iteration "learnGoWithTest/internal/iteration"
)

type Employee struct {
	Name string
	Age int
	Salary float64
}

func main() {
	fmt.Println("Learn Go with TDD [Test Driven Development]")

	fmt.Println(hello.Hello("world !", ""))
	// Improve this and customize it
	fmt.Println(hello.CustomPrint("Test this function as it is .... "))

	fmt.Println("String function implemented with Iteration logic")
	fmt.Printf("Repeation of character a 5 times: %q \n", iteration.Repeat("a", 5))
	fmt.Println(strings.Count("cheese", "e"))
	fmt.Println(strings.Count("five", ""))
	fmt.Println(iteration.CharacterCount("mayanamara", "t"))

	arr := []int{1, 2, 3, 4, 5}
	fmt.Printf("Get the sum of array %d \n", array.Sum(arr))
	arr1 := []int{1, 2, 3}
	arr2 := []int{4, 5, 6}

	fmt.Printf("Go the sum of two array %v \n. ", array.SumAll(arr1, arr2))

	employee1 := Employee{
		Name: "amolpratap singh",
		Age: 29,
		Salary: 100000000,
	}

	employee2 := Employee{
		Name: "rajesh singh",
		Age: 56,
		Salary: 123456889.01,
	}

	fmt.Printf("Employee 1 detail : %v", employee1)
	fmt.Printf("Employee 2 detail : %v", employee2)
}
