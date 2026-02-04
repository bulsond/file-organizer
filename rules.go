package main

// DefaultRules правила для сортировки файлов
type DefaultRules map[string]string

// NewDefaultRules создание экземпляра правил сортировки
func NewDefaultRules() DefaultRules {
	return DefaultRules{
		// Music
		".wav": "Music",
		".mp3": "Music",

		// Video
		".mp4": "Video",
		".avi": "Video",

		// Archives
		".rar": "Archives",
		".zip": "Archives",

		// Images
		".jpg":  "Images",
		".jpeg": "Images",
		".png":  "Images",

		// Documents
		".pdf":  "Documents",
		".doc":  "Documents",
		".docx": "Documents",
		".txt":  "Documents",
	}
}
