package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
		dir := input
		if len(dir) == 0 {
			dir = a.currentDir
		}
		path, err := resolvePath(dir)
		if err != nil {
			fmt.Printf("Ошибка разрешения пути: %v\n", err)
			continue
		}

		fmt.Printf("Выбран каталог %s\n", path)
		// сортируем
		organizeDir(path)
	}
	fmt.Println("Успешных договоров о межгалактическом сотрудничестве!")
}

// resolvePath преобразует путь в абсолютный
func resolvePath(dir string) (string, error) {
	// если уже абсолютный путь
	if filepath.IsAbs(dir) {
		return filepath.Clean(dir), nil
	}

	// начинается с ~ (Unix домашняя директория)
	if strings.HasPrefix(dir, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		// замена тильды
		dir = filepath.Join(home, dir[1:])
	}

	// преобразуем относительный в абсолютный
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	return filepath.Clean(absPath), nil
}

// organizeDir сортировать файлы в каталоге
func organizeDir(path string) {
	fo, err := NewFileOrganizer(path)
	if err != nil {
		fmt.Printf("Не удалось начать сортировку файлов по причине: %v\n",
			err)
	}
	defer fo.Close()

	// работаем
	fmt.Println("Начинаю сортировку файлов...")
	if err := fo.Organize(); err != nil {
		fmt.Printf("Ошибка при сортировке: %v\n", err)
	} else {
		fmt.Println("Сортировка завершена успешно!")
	}

	fo.Report()
}
