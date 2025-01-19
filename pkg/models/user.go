package models

import(
	"time"
	"encoding/json"
	jwt "github.com/golang-jwt/jwt"
)

type User struct {
	ID    		int    				`json:"id"`
	Name  		string 				`json:"name"`						// ФИО
	Email 		string 				`json:"email"`						// Email
	Phone 		*string 			`json:"phone"`						// Телефон
	Password 	string 				`json:"password"`					// Пароль
	City        *string 	   	   	`json:"city"`           			// Город
	Logo       	*string    			`json:"logo"`                 		// Логотип
	Media       *json.RawMessage    `json:"media"`                		// Медиа
	Status 		int 				`json:"status"`						// Статус
	CreatedAt   time.Time       	`json:"created_at"`              	// Дата создания
	UpdatedAt   time.Time       	`json:"updated_at"`             	// Дата последнего обновления
	DeletedAt 	*time.Time 	    	`json:"deleted_at, omitempty"`		// Дата удаления
}

type UserView struct {
	ID    		int    				`json:"id"`
	Name  		string 				`json:"name"`						// ФИО
	Email 		string 				`json:"email"`						// Email
	Phone 		*string 			`json:"phone"`						// Телефон
	City        *string 	   	   	`json:"city"`           			// Город
	Logo       	*string    			`json:"logo"`                 		// Логотип
	Media       *json.RawMessage    `json:"media"`                		// Медиа
	Status 		int 				`json:"status"`						// Статус
	CreatedAt   time.Time       	`json:"created_at"`              	// Дата создания
}

type CreateUserRequest struct {
	ID    		int
	Name 		string 				`json:"name" validate:"required,min=3,max=128"`
	Email 		string 				`json:"email" validate:"required,email"`
	Phone 		string 				`json:"phone, omitempty" validate:"min=6,max=20"`
	Password	string 				`json:"password" validate:"required,min=4,max=128"`
}

type UpdateUserRequest struct {
	ID    		int    				`json:"id"`
	Name 		string 				`json:"name" validate:"required,min=3,max=128"`
	Email 		string 				`json:"email" validate:"required,email"`
	Phone 		*string 	 		`json:"phone" validate:"min=6,max=20"`
	Password	string 				`json:"password" validate:"required,min=4,max=128"`
}


type LoginUserRequest struct {
	Email 		string 				`json:"email" validate:"required,email"`
	Password	string 				`json:"password" validate:"required"`
}

// Claims содержит информацию, которую мы хотим включить в токен
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}