package main

import(
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone *string `json:"phone"`
	Status int `json:"status"`
	DataCreate time.Time `json:"created_at"`
	DateUpdate time.Time `json:"updated_at"`
	DateDelete *time.Time `json:"deleted_at"`
}
