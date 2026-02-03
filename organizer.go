package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
	statistics     map[string]*FileStats
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
		statistics:     map[string]*FileStats{},
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

// moveFile перемещение файла, возвращает размер перемещенного файла
func (fo *FileOrganizer) moveFile(sourcePath, targetDir string) (int64, error) {
	var msg string
	var size int64
	// существование исходного файла
	origInfo, err := os.Stat(sourcePath)
	if os.IsNotExist(err) {
		msg = fmt.Sprintf("Исходный файл не найден: %s", sourcePath)
		fo.logError(msg)
		return size, errors.New(msg)
	}
	size = origInfo.Size()

	// целевая директория
	fullTargetDir := filepath.Join(fo.sourceDir, targetDir)
	if err := os.MkdirAll(fullTargetDir, 0775); err != nil {
		msg = fmt.Sprintf("Не удалось создать директорию %s: %v", targetDir, err)
		fo.logError(msg)
		return size, errors.New(msg)
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
		fileName = generateFileName(fileName)
		targetPath = filepath.Join(fullTargetDir, fileName)
	}

	// перемещение файла
	if err := os.Rename(sourcePath, targetPath); err != nil {
		msg = fmt.Sprintf("Не удалось переместить файл %s: %v", fileName, err)
		fo.logError(msg)
		return size, errors.New(msg)
	}

	msg = fmt.Sprintf("Результат: %s/%s", targetDir, fileName)
	fo.logSuccess(msg)

	return size, nil
}

// Organize сортировка файлов в целевой дериктории
func (fo *FileOrganizer) Organize() error {
	fo.logSuccess("Начинаем сортировку файлов в директории: " + fo.sourceDir)

	err := filepath.WalkDir(fo.sourceDir, fo.walkWoker)
	if err != nil {
		fo.logError(fmt.Sprintf("Ошибка при сканировании директории: %v", err))
		return fmt.Errorf("ошибка сканирования: %w", err)
	}

	fo.logSuccess("Завершили сортировку файлов")
	return nil
}

// walkWoker обработчик для filepath.WalkDir
func (fo *FileOrganizer) walkWoker(path string, d fs.DirEntry, err error) error {
	var msg string
	// Пропускаем ошибки доступа (например, к недоступным директориям)
	if err != nil {
		msg = fmt.Sprintf("Ошибка доступа к %s: %v", path, err)
		fo.logError(msg)
		return nil
	}

	// Пропускаем директории — нам нужны только файлы
	if path == fo.sourceDir {
		return nil
	}
	if d.IsDir() {
		return fs.SkipDir // пропускаем содержимое поддиректорий
	}

	// Проверяем расширение файла
	ext := filepath.Ext(path)
	if len(ext) == 0 {
		msg = fmt.Sprintf("Файл не имеет расширения: %s",
			filepath.Base(path))
		fo.logError(msg)
		return nil
	}

	targetDir, exists := fo.rulesMap[ext]
	if !exists {
		msg = fmt.Sprintf("Неизвестное расширение %s у файла: %s",
			ext, filepath.Base(path))
		fo.logError(msg)
		return nil
	}

	// Перемещаем файл
	_, errM := fo.moveFile(path, targetDir)
	if errM != nil {
		msg = fmt.Sprintf("Не удалось переместить файл %s: %v",
			filepath.Base(path), err)
		fo.logError(msg)
		return nil
	}

	// увеличиваем счетчик файлов
	fo.processedFiles++

	return nil // продолжаем обход
}

func (fo *FileOrganizer) setStats() {

}
