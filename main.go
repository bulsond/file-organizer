package main

import "fmt"

func main() {
	const dir = "generated_files"
	const dubDir = "."
	fmt.Println("Привет")

	fo, err := NewFileOrganizer(dubDir)
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	fmt.Println(fo)

	// работаем
	// fmt.Println("Начинаю сортировку файлов...")
	// if err := fo.Organize(); err != nil {
	// 	fmt.Printf("Ошибка при сортировке: %v\n", err)
	// } else {
	// 	fmt.Println("Сортировка завершена успешно!")
	// }
}
