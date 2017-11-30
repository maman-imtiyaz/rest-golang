package main

import "fmt"

type person struct {
	Name string
	Age  int
}

func main() {
	fmt.Println("1 >>", person{"maman", 32})
	fmt.Println("2 >>", person{Name:"maman"})
	fmt.Println("3 >>", person{Age:32})

	temp := person{"maman",32}
	fmt.Println("4 >>", temp.Name)

	temp2 := temp
	temp2.Name = "vry"
	fmt.Println("5 >>", temp2.Name, " | ", temp2.Age)
}