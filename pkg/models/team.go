package models

import(
	"time"
	"encoding/json"
)

// Team - структура для описания команды
type Team struct {
	ID              	int64      	   		`json:"id"`
	Name            	string         		`json:"name"`           			// Название
	Description     	*string 	   		`json:"description"`          		// Описание
	City            	string 	   	   		`json:"city"`           			// Город
	UniformColor    	*string 	   		`json:"uniform_color"`           	// Цвет формы
	ParticipantCount 	*int  		   		`json:"participant_count"`    		// Кол-во участников
	Responsible     	*int          		`json:"responsible"`          		// Ответственный
	DisabilityCategory 	*string				`json:"disability_category"`  		// Категория инвалидности
	Logo            	*json.RawMessage    `json:"logo"`                 		// Логотип
	Media           	*json.RawMessage    `json:"media"`                		// Медиа
	Status 				*int 		   		`json:"status"`						// Статус
	CreatedAt       	time.Time      		`json:"created_at"`              	// Дата создания
	UpdatedAt       	time.Time      		`json:"updated_at"`             	// Дата последнего обновления
	DeletedAt 			*time.Time 	   		`json:"deleted_at, omitempty"`		// Дата удаления
}

// CreateTeamRequest - структура запроса на добавление команды
type CreateTeamRequest struct {
	Name            	string         		`json:"name" validate:"required"`       // Название*
	Description     	*string 	   		`json:"description, omitempty"`          // Описание
	City            	string         		`json:"city" validate:"required"`        // Город
	UniformColor    	*string 	   		`json:"uniform_color"`           		// Цвет формы
	ParticipantCount 	*int  		   		`json:"participant_count"`       		// Кол-во участников
	Responsible     	*int          		`json:"responsible"`             			// Ответственный
	DisabilityCategory  *string 	   		`json:"disability_category"`     		// Категория инвалидности
	Logo            	*json.RawMessage    `json:"logo"`                  // Логотип
	Media           	*json.RawMessage    `json:"media"`                   // Медиа
}

// UpdateTeamRequest - структура запроса на изменение команды
type UpdateTeamRequest struct {
	Name            	string         		`json:"name" validate:"required"`       // Название*
	Description     	*string 	   		`json:"description, omitempty"`          // Описание
	City            	string         		`json:"city" validate:"required"`        // Город
	UniformColor    	*string 	   		`json:"uniform_color"`           		// Цвет формы
	ParticipantCount 	*int  		   		`json:"participant_count"`       		// Кол-во участников
	Responsible     	*int          		`json:"responsible"`             		// Ответственный
	DisabilityCategory  *string 	   		`json:"disability_category"`     		// Категория инвалидности
	Logo            	*json.RawMessage    `json:"logo"`                    		// Логотип
	Media           	*json.RawMessage    `json:"media"`                   		// Медиа
}

