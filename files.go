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

// generateFileName создание нового имени файла с постфиксом
func generateFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)
	postfix := time.Now().Format("2006-03-15_15-04-05")

	return fmt.Sprintf("%s_%s%s", name, postfix, ext)
}
