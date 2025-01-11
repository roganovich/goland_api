package models

import(
	"time"
)

// Media - структура для медиа-файлов
type Media struct {
	ID      	int     	`json:"id"`
	Name      	string    	`json:"name"`
	Path      	string     	`json:"path"`
	Ext 		string     	`json:"ext"`
	Size     	int64      	`json:"size"`
	CreatedAt 	time.Time 	`json:"created_at"`
}