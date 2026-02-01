package main

import (
	"errors"
	"os"
	"path/filepath"
)

// FileOrganizer базовая структура органайзера
type FileOrganizer struct {
	sourceDir      string
	rulesMap       map[string]string
	processedFiles int
	logFile        *os.File
}

// NewFileOrganizer создание FileOrganizer
func NewFileOrganizer(sourceDir string) (*FileOrganizer, error) {
	const nameLogFile = "organizer.log"
	if len(sourceDir) == 0 {
		return &FileOrganizer{},
			errors.New("Не указан путь к директории с файлами для сортировки")
	}

	// существование директории
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return &FileOrganizer{}, err
	}

	// абсолютный путь к директории
	absPath, err := filepath.Abs(sourceDir)
	if err != nil {
		return &FileOrganizer{}, err
	}

	// является ли директорией
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return &FileOrganizer{}, err
	}
	if !fileInfo.IsDir() {
		return &FileOrganizer{},
			errors.New("Не является путем к директории с файлами для сортировки ")
	}

	// лог файл
	path := filepath.Join(absPath, nameLogFile)
	logFile, err := NewLogFile(path)
	if err != nil {
		return &FileOrganizer{}, err
	}

	return &FileOrganizer{
		sourceDir:      absPath,
		rulesMap:       NewDefaultRules(),
		processedFiles: 0,
		logFile:        logFile,
	}, nil
}

// Close закрытие структуры
func (fo *FileOrganizer) Close() error {
	if fo.logFile != nil {
		return fo.logFile.Close()
	}
	return nil
}

// NewLogFile создание ссылки на лог файл
func NewLogFile(filePath string) (*os.File, error) {
	if len(filePath) == 0 {
		return nil, errors.New("Не указан путь к файлу логов")
	}

	logFile, err := os.OpenFile(
		filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	return logFile, nil
}
