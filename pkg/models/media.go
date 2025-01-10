package models

import(
	"time"
)

// Media - структура для медиа-файлов
type Media struct {
	URL      	string     	`json:"url"`
	FileType 	string     	`json:"file_type"`
	Size     	int64      	`json:"size"`
	CreatedAt 	time.Time 	`json:"created_at"`
}