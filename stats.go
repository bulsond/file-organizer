package main

import "fmt"

// FileStats для сбора статистики по директориям
type FileStats struct {
	Count     int
	TotalSize int64
}

// String() реализация fmt.Stringer
func (fs FileStats) String() string {
	return fmt.Sprintf("  - Количество файлов: %d\n  - Общий размер: %s",
		fs.Count,
		stringTotalSize(fs.TotalSize),
	)
}

// stringTotalSize строковое отображение общего размера файлов
func stringTotalSize(size int64) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%d B", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
	default:
		return fmt.Sprintf("%.1f GB", float64(size)/(1024*1024*1024))
	}
}
