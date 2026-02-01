package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	SuccessMsg = "[SUCCESS]"
	ErrorMsg   = "[ERROR]"
	SystemMsg  = "[SYSTEM]"
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
		// Создаем логгер для записи завершения сессии
		logger := log.New(fo.logFile, "", 0)
		logger.Printf("%s Сессия завершена", SystemMsg)

		return fo.logFile.Close()
	}
	return nil
}

// NewLogFile создание ссылки на лог файл
func NewLogFile(logPath string) (*os.File, error) {
	if len(logPath) == 0 {
		return nil, errors.New("Не указан путь к файлу логов")
	}

	logFile, err := os.OpenFile(
		logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	// Создаем логгер для записи начала сессии
	logger := log.New(logFile, "", 0)
	logger.Printf("%s Сессия начата", SystemMsg)

	return logFile, nil
}

// logSuccess создать запись об успешной операции
func (fo *FileOrganizer) logSuccess(message string) {
	fo.writeRecord(SuccessMsg, message)
}

// logError создать запись об ошибке
func (fo *FileOrganizer) logError(message string) {
	fo.writeRecord(ErrorMsg, message)
}

// writeRecord создать запись в лог
func (fo *FileOrganizer) writeRecord(msgType, message string) {
	var logger *log.Logger
	if fo.logFile == nil {
		logger = log.Default()
	}
	// Создаем логгер с временными метками
	logger = log.New(fo.logFile, "", log.LstdFlags)
	// Логгер автоматически добавит временную метку
	logger.Printf("%s %s", msgType, message)
}

// moveFile перемещение файла
func (fo *FileOrganizer) moveFile(sourcePath, targetDir string) error {
	var msg string
	// существование исходного файла
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		msg = fmt.Sprintf("Исходный файл не найден: %s", sourcePath)
		fo.logError(msg)
		return errors.New(msg)
	}

	// целевая директория
	fullTargetDir := filepath.Join(fo.sourceDir, targetDir)
	if err := os.MkdirAll(fullTargetDir, 0775); err != nil {
		msg = fmt.Sprintf("Не удалось создать директорию %s: %v", targetDir, err)
		fo.logError(msg)
		return errors.New(msg)
	}
	msg = fmt.Sprintf("Целевая директория: %s", targetDir)
	fo.logSuccess(msg)

	// имя файла
	fileName := filepath.Base(sourcePath)
	msg = fmt.Sprintf("Исходный файл: %s", fileName)
	fo.logSuccess(msg)

	// полный путь к целевому файлу
	targetPath := filepath.Join(fullTargetDir, fileName)

	// есть ли конфликт имен
	if _, err := os.Stat(targetPath); err == nil {
		// значит такой уже есть, создаем новое имя
		msg = fmt.Sprintf("Существующий файл: %s/%s", targetDir, fileName)
		fo.logSuccess(msg)
		fileName = fo.generateFileName(fileName)
		targetPath = filepath.Join(fullTargetDir, fileName)
	}

	// перемещение файла
	if err := os.Rename(sourcePath, targetPath); err != nil {
		msg = fmt.Sprintf("Не удалось переместить файл %s: %v", fileName, err)
		fo.logError(msg)
		return errors.New(msg)
	}

	msg = fmt.Sprintf("Результат: %s/%s", targetDir, fileName)
	fo.logSuccess(msg)

	return nil
}

// generateFileName создание нового имени файла с постфиксом
func (fo *FileOrganizer) generateFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)
	postfix := time.Now().Format("2006-03-15_15-04-05")

	return fmt.Sprintf("%s_%s%s", name, postfix, ext)
}
