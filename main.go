package main

import "fmt"

func main() {
	fmt.Println("Привет")

	rules := NewDefaultRules()
	for ext, folder := range rules {
		fmt.Printf("Расширение: %s -> Папка: %s\n", ext, folder)
	}

	fmt.Println(NewFileOrganizer("."))
}
