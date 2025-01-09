package models

import(
	"time"
)

type User struct {
	ID    		int    			`json:"id"`
	Name  		string 			`json:"name"`
	Email 		string 			`json:"email"`
	Phone 		*string 		`json:"phone"`
	Status 		int 			`json:"status"`
	CreatedAt   time.Time       `json:"created_at"`              	// Дата создания
	UpdatedAt   time.Time       `json:"updated_at"`             	// Дата последнего обновления
	DeletedAt 	*time.Time 	    `json:"deleted_at, omitempty"`		// Дата удаления
}

type CreateUserRequest struct {
	Name 	string 				`json:"name" validate:"required,min=3,max=128"`
	Email 	string 				`json:"email" validate:"required,email"`
	Phone 	*string 			`json:"phone" validate:"min=6,max=20"`
}

type UpdateUserRequest struct {
	Name 	string 				`json:"name" validate:"required,min=3,max=128"`
	Email 	string 				`json:"email" validate:"required,email"`
	Phone 	*string 	 		`json:"phone" validate:"min=6,max=20"`
}
