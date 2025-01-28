package models

import(
	"database/sql"
	"time"
	"encoding/json"
)

type Field struct {
	ID          int64				`json:"id"`
	Name        string       		`json:"name"`
	Description sql.NullString   	`json:"description"`
	City        string       		`json:"city"`
	Address     string       		`json:"address"`
	Location    []byte       		`json:"location"`
	Square      sql.NullInt64 		`json:"square"`
	Info        sql.NullString 		`json:"info"`
	Places      int           		`json:"places"`
	Dressing    bool          		`json:"dressing"`
	Toilet      bool          		`json:"toilet"`
	Display     bool          		`json:"display"`
	Parking     bool          		`json:"parking"`
	ForDisabled bool          		`json:"for_disabled"`
	Logo        sql.NullString 		`json:"logo"`
	Media       []byte        		`json:"media"`
	Status      int16         		`json:"status"`
	CreatedAt   time.Time     		`json:"created_at"`
	UpdatedAt   time.Time     		`json:"updated_at"`
	DeletedAt   *time.Time    		`json:"deleted_at"`
}

type FieldView struct {
	ID          int64				`json:"id"`
	Name        string       		`json:"name"`
	Description sql.NullString   	`json:"description"`
	City        string       		`json:"city"`
	Address     string       		`json:"address"`
	Location    []byte       		`json:"location"`
	Square      sql.NullInt64 		`json:"square"`
	Info        sql.NullString 		`json:"info"`
	Places      int           		`json:"places"`
	Dressing    bool          		`json:"dressing"`
	Toilet      bool          		`json:"toilet"`
	Display     bool          		`json:"display"`
	Parking     bool          		`json:"parking"`
	ForDisabled bool          		`json:"for_disabled"`
	Logo        sql.NullString 		`json:"logo"`
	Media       []byte        		`json:"media"`
	Status      int16         		`json:"status"`
	CreatedAt   time.Time     		`json:"created_at"`
}

type CreateFieldRequest struct {
	Name        string       		`json:"name" validate:"required"`
	Description sql.NullString   	`json:"description"`
	City        string       		`json:"city" validate:"required"`
	Address     string       		`json:"address" validate:"required"`
	Location    []byte       		`json:"location"`
	Square      sql.NullInt64 		`json:"square"`
	Info        sql.NullString 		`json:"info"`
	Places      int           		`json:"places"`
	Dressing    bool          		`json:"dressing"`
	Toilet      bool          		`json:"toilet"`
	Display     bool          		`json:"display"`
	Parking     bool          		`json:"parking"`
	ForDisabled bool          		`json:"for_disabled"`
	Logo        sql.NullString 		`json:"logo"`
	Media       []byte        		`json:"media"`
}

type UpdateFieldRequest struct {
	Name        string       		`json:"name" validate:"required"`
	Description sql.NullString   	`json:"description"`
	City        string       		`json:"city" validate:"required"`
	Address     string       		`json:"address" validate:"required"`
	Location    []byte       		`json:"location"`
	Square      sql.NullInt64 		`json:"square"`
	Info        sql.NullString 		`json:"info"`
	Places      int           		`json:"places"`
	Dressing    bool          		`json:"dressing"`
	Toilet      bool          		`json:"toilet"`
	Display     bool          		`json:"display"`
	Parking     bool          		`json:"parking"`
	ForDisabled bool          		`json:"for_disabled"`
	Logo        sql.NullString 		`json:"logo"`
	Media       []byte        		`json:"media"`
}