package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type App struct {
	currentDir  string
	reader      *bufio.Reader
	quitCommand string
}

// NewApp создание экземляра приложения
func NewApp() (*App, error) {
	// абсолютный путь текущей директории
	dir, err := os.Getwd()
	if err != nil {
		return &App{}, err
	}
	return &App{
		currentDir:  dir,
		reader:      bufio.NewReader(os.Stdin),
		quitCommand: "quit",
	}, nil
}

// Run запуск приложения
func (a *App) Run() {
	fmt.Println("ПОРТ")
	fmt.Println("Программа Организации Разнообразных Типов")
	fmt.Println()
	fmt.Printf("Текущий каталог: %s\n", a.currentDir)
	fmt.Printf("Введите '%s' для выхода\n", a.quitCommand)

	for {
		fmt.Print("Введите путь к каталогу (Enter для текущего): ")

		// читаем ввод пользователя
		input, err := a.reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Ошибка чтения ввода: %v\n", err)
			continue
		}

		// проверяем ввод
		input = strings.TrimSpace(input)
		if strings.EqualFold(input, a.quitCommand) {
			fmt.Println("Выход из программы...")
			break
		}

		// целевой каталог
		targetDir := input
		if len(targetDir) == 0 {
			targetDir = a.currentDir
		}
		fmt.Printf("Выбран каталог %s\n", targetDir)
	}
	fmt.Println("Успешных договоров о межгалактическом сотрудничестве!")
}
