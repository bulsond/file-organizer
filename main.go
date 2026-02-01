package main

import "fmt"

func main() {
	fmt.Println("Привет")

	rules := NewDefaultRules()
	for ext, folder := range rules {
		fmt.Printf("Расширение: %s -> Папка: %s\n", ext, folder)
	}

	fileOrganizer, err := NewFileOrganizer(".")
	if err != nil {
		panic(err)
	}
	defer fileOrganizer.Close()
	fmt.Println(fileOrganizer)

	fileOrganizer.logSuccess("Выполнена успешная операция")
	fileOrganizer.logError("Выполнение вызвало ошибку")

}
