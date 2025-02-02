package handlers

import (
	"goland_api/pkg/models"
	"goland_api/pkg/database"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"database/sql"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Документация для метода GetFields
// @Summary Возвращает список всех команд
// @Description Получение списка всех команд
// @Tags Команды
// @Accept application/json
// @Produces application/json
// @Success 200 {object} []models.FieldView
// @Failure 400 Bad Request
// @Failure 500 Internal Server Error
// @Router /api/fields [get]
func GetFields() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := database.DB.Query("SELECT id, name, description, city, address, logo, media, status, created_at  FROM fields")
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		fields := []models.FieldView{}
		for rows.Next() {
			var fieldView models.FieldView
			var logo sql.NullString
			var media json.RawMessage

			if err := rows.Scan(
				&fieldView.ID,
				&fieldView.Name,
				&fieldView.Description,
				&fieldView.City,
				&fieldView.Address,
				&logo,
				&media,
				&fieldView.Status,
				&fieldView.CreatedAt,
				); err != nil {
				log.Println(err)
			}

			if (logo.Valid){
				var logoFile models.Media
				errorMedia, logoFile := getOneMedia(logo.String)
				if  errorMedia != nil {
					log.Println(errorMedia.Error())
				}else{
					fieldView.Logo = &logoFile
				}
			}
			if (media != nil && len(media) > 0){
				var mediaList []models.Media
				var mediaFiles []string
				err := json.Unmarshal(media, &mediaFiles)
				if err != nil {
					log.Println("Ошибка при парсинге JSON:", err)
				}
				for _, mediaFile := range mediaFiles {
					errorMedia, mediaFile := getOneMedia(mediaFile)
					if  errorMedia != nil {
						log.Println(errorMedia.Error())
					}else{
						mediaList = append(mediaList, mediaFile)
					}
				}
				fieldView.Media = &mediaList
			}
			fields = append(fields, fieldView)
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
		}

		json.NewEncoder(w).Encode(fields)
	}
}

func getOneFieldById(paramId int64) (error, models.FieldView) {
	var fieldView models.FieldView
	var logo sql.NullString
	var media json.RawMessage

	err := database.DB.QueryRow("SELECT id, name, description, city, address, logo, media, status, created_at FROM fields WHERE id = $1", int64(paramId)).Scan(
		&fieldView.ID,
		&fieldView.Name,
		&fieldView.Description,
		&fieldView.City,
		&fieldView.Address,
		&logo,
		&media,
		&fieldView.Status,
		&fieldView.CreatedAt,
	)
	if err != nil {
		return err, fieldView
	}

	if (logo.Valid){
		var logoFile models.Media
		errorMedia, logoFile := getOneMedia(logo.String)
		if  errorMedia != nil {
			log.Println(errorMedia.Error())
		}else{
			fieldView.Logo = &logoFile
		}
	}
	if (media != nil && len(media) > 0){
		var mediaList []models.Media
		var mediaFiles []string
		err := json.Unmarshal(media, &mediaFiles)
		if err != nil {
			log.Println("Ошибка при парсинге JSON:", err)
		}
		for _, mediaFile := range mediaFiles {
			errorMedia, mediaFile := getOneMedia(mediaFile)
			if  errorMedia != nil {
				log.Println(errorMedia.Error())
			}else{
				mediaList = append(mediaList, mediaFile)
			}
		}
		fieldView.Media = &mediaList
	}

	return err, fieldView
}

// Документация для метода GetField
// @Summary Возвращает информацию о команде по ID
// @Description Получение информации о команде по идентификатору
// @Tags Команды
// @Param id path int true "ID команды"
// @Success 200 {object} models.FieldView
// @Failure 400 Bad Request
// @Failure 404 Not Found
// @Router /api/fields/{id} [get]
func GetField() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])

		errorResponse, fieldView := getOneFieldById(int64(paramId))
		if  errorResponse != nil {
			http.Error(w, errorResponse.Error(), http.StatusBadRequest)
			return
		}


		json.NewEncoder(w).Encode(fieldView)
	}
}

func validateCreateFieldRequest(r *http.Request) (error, models.CreateFieldRequest) {
	var req models.CreateFieldRequest
	if validation := json.NewDecoder(r.Body).Decode(&req); validation != nil {
		return validation, req
	}
	validate := validator.New()
	if validation := validate.Struct(req); validation != nil {
		return validation, req
	}

	return nil, req
}

func validateUpdatedAtFieldRequest(r *http.Request) (error, models.UpdateFieldRequest) {
	var req models.UpdateFieldRequest
	if validation := json.NewDecoder(r.Body).Decode(&req); validation != nil {
		return validation, req
	}
	validate := validator.New()
	if validation := validate.Struct(req); validation != nil {
		return validation, req
	}

	return nil, req
}

// Документация для метода CreateField
// @Summary Создание новой команды
// @Description Создание новой команды
// @Tags Команды
// @Param createField body models.CreateFieldRequest true "Данные для создания новой команды"
// @Consumes application/json
// @Produces application/json
// @Success 201 {object} models.FieldView
// @Failure 422 Unprocessable Entity
// @Router /api/fields [post]
func CreateField() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validation, fieldRequest := validateCreateFieldRequest(r)
		if  validation != nil {
			http.Error(w, validation.Error(), http.StatusBadRequest)
			return
		}

		var field models.Field
		field.Name = fieldRequest.Name
		field.Description = fieldRequest.Description
		field.City = fieldRequest.City
		field.Address = fieldRequest.Address
		field.Logo = fieldRequest.Logo
		field.Media = fieldRequest.Media
		err := database.DB.QueryRow("INSERT INTO fields (name, description, city, address, logo, media) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			field.Name,
			field.Description,
			field.City,
			field.Address,
			field.Logo,
			field.Media,
		).Scan(&field.ID)
		if err != nil {
			log.Println(err)
		}

		errField, fieldView := getOneFieldById(int64(field.ID))
		if errField != nil {
			http.Error(w, errField.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(fieldView)
	}
}

// Документация для метода UpdateField
// @Summary Обновление существующей команды
// @Description Обновление существующей команды
// @Tags Команды
// @Param updateField body models.UpdateFieldRequest true "Данные для обновления команды"
// @Consumes application/json
// @Produces application/json
// @Param id path int true "ID команды"
// @Success 204 No Content
// @Failure 422 Unprocessable Entity
// @Failure 404 Not Found
// @Router /api/fields/{id} [put]
func UpdateField() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validation, fieldRequest := validateUpdatedAtFieldRequest(r)
		if  validation != nil {
			http.Error(w, validation.Error(), http.StatusBadRequest)
			return
		}
		var field models.Field
		field.Name = fieldRequest.Name
		field.Description = fieldRequest.Description
		field.City = fieldRequest.City
		field.Address = fieldRequest.Address
		field.Logo = fieldRequest.Logo
		field.Media = fieldRequest.Media
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])
		field.ID = int64(paramId)

		_, errUpdate := database.DB.Exec("UPDATE fields SET name = $1, description = $2, city = $3, address = $4, logo = $5, media = $6 WHERE id = $7",
			field.Name,
			field.Description,
			field.City,
			field.Address,
			field.Logo,
			field.Media,
			paramId)
		if errUpdate != nil {
			log.Println(errUpdate)
			http.Error(w, errUpdate.Error(), http.StatusBadRequest)
		}

		errorResponse, fieldView := getOneFieldById(int64(paramId))
		if  errorResponse != nil {
			http.Error(w, errorResponse.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(fieldView)
	}
}

// Документация для метода DeleteField
// @Summary Удаляет команду по ID
// @Description Удаление команды по идентификатору
// @Tags Команды
// @Param id path int true "ID команды"
// @Success 204 No Content
// @Failure 404 Not Found
// @Router /api/fields/{id} [delete]
func DeleteField() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])
		errorResponse, fieldView := getOneFieldById(int64(paramId))
		if errorResponse != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := database.DB.Exec("DELETE FROM fields WHERE id = $1", fieldView.ID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode("Field deleted")
		}
	}
}


