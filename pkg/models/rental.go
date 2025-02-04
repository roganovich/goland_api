package models

import(
	"time"
)

// Определяем структуру "Аренда"
type Rental struct {
	ID    		 int64   		`json:"id"`			   // Идентификатор
	FieldID      int64       	`json:"field_id"`      // Идентификатор площадки
	TeamID       int64       	`json:"team_id"`       // Идентификатор команды
	UserID       int64       	`json:"user_id"`       // Идентификатор пользователя
	Comment      string    		`json:"comment"`       // Комментарий
	StartDate    time.Time 		`json:"start_date"`    // Дата начала аренды
	EndDate      time.Time 		`json:"end_date"`      // Дата завершения аренды
	Duration     int       		`json:"duration"`      // Длительность аренды (например, в часах)
	Status       int    		`json:"status"`        // Статус аренды
	CreatedAt    time.Time      `json:"created_at"`    // Дата создания
	UpdatedAt    time.Time      `json:"updated_at"`    // Дата последнего обновления
	DeletedAt 	 *time.Time 	`json:"deleted_at"`	   // Дата удаления
}

type RentalView struct {
	ID    		 int64   		`json:"id"`			   // Идентификатор
	Field     	 FieldView      `json:"field"`         // Пощадка
	Team       	 TeamView       `json:"team"`          // Команда
	User      	 UserView       `json:"user"`          // Пользователь
	Comment      string    		`json:"comment"`       // Комментарий
	StartDate    time.Time 		`json:"start_date"`    // Дата начала аренды
	EndDate      time.Time 		`json:"end_date"`      // Дата завершения аренды
	Duration     int       		`json:"duration"`      // Длительность аренды (например, в часах)
	Status       int    		`json:"status"`        // Статус аренды
	CreatedAt    time.Time      `json:"created_at"`    // Дата создания
}

type CreateRentalRequest struct {
	ID    		 int64   		`json:"id"`			   					   // Идентификатор
	FieldID      int64       	`json:"field_id" validate:"required"`      // Идентификатор площадки
	TeamID       int64       	`json:"team_id" validate:"required"`       // Идентификатор команды
	Comment      string    		`json:"comment"`       					   // Статус аренды
	StartDate    string 		`json:"start_date" validate:"required"`    // Дата начала аренды
	EndDate      string 		`json:"end_date" validate:"required"`      // Дата завершения аренды
}

type RentalSearchFilter struct {
	Search 		*string   	`json:"search"` 		// Поисковый запрос
	FieldId   	*[]int64 	`json:"field_ids"`   	// Фильтр по площадкам
	TeamId   	*[]int64 	`json:"team_ids"`   	// Фильтр по командам
	StartDate   *string 	`json:"start_date"`     // Дата начала аренды
	EndDate     *string 	`json:"end_date"`       // Дата завершения аренды
	Status      *int    	`json:"status"`         // Статус аренды
}