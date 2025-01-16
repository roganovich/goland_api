package handlers

import (
	"goland_api/pkg/models"

	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Документация для метода GetUsers
// @Summary Возвращает список всех пользователей
// @Description Получение списка всех пользователей
// @Tags Пользователи
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} []models.User
// @Failure 400 Bad Request
// @Failure 500 Internal Server Error
// @Router /api/users [get]
func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		users := []models.User{}
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(users)
	}
}

func getOneUser(db *sql.DB, paramId int) (error, models.User) {
	var user models.User
	err := db.QueryRow("SELECT * FROM users WHERE id = $1", int64(paramId)).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		)

	return err, user
}


// Документация для метода GetUser
// @Summary Возвращает информацию о пользователе по ID
// @Description Получение информации о пользователе по идентификатору
// @Tags Пользователи
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.User
// @Failure 400 Bad Request
// @Failure 404 Not Found
// @Router /api/users/{id} [get]
func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])

		errorResponse, user := getOneUser(db, paramId)
		if  errorResponse != nil {
			http.Error(w, errorResponse.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func validateCreateUserRequest(r *http.Request) (error, models.CreateUserRequest) {
	var req models.CreateUserRequest
	if validation := json.NewDecoder(r.Body).Decode(&req); validation != nil {
		return validation, req
	}
	validate := validator.New()
	if validation := validate.Struct(req); validation != nil {
		return validation, req
	}

	return nil, req
}

func valiUpdatedAtUserRequest(r *http.Request) (error, models.UpdateUserRequest) {
	var req models.UpdateUserRequest
	if validation := json.NewDecoder(r.Body).Decode(&req); validation != nil {
		return validation, req
	}
	validate := validator.New()
	if validation := validate.Struct(req); validation != nil {
		return validation, req
	}

	return nil, req
}

// Документация для метода CreateUser
// @Summary Создание нового пользователя
// @Description Создание нового пользователя
// @Tags Пользователи
// @Param createUser body models.CreateUserRequest true "Данные для создания пользователя"
// @Consumes application/json
// @Produces application/json
// @Success 201 {object} models.User
// @Failure 422 Unprocessable Entity
// @Router /api/users [post]
func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validation, userRequest := validateCreateUserRequest(r)
		if  validation != nil {
			http.Error(w, validation.Error(), http.StatusBadRequest)
			return
		}
		var user models.User
		user.Name = userRequest.Name
		user.Email = userRequest.Email
		user.Phone = userRequest.Phone

		err := db.QueryRow("INSERT INTO users (name, email, phone) VALUES ($1, $2, $3) RETURNING id", user.Name, user.Email, user.Phone).Scan(&user.ID)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(user)
	}
}

// Документация для метода UpdateUser
// @Summary Обновление существующего пользователя
// @Description Обновление существующего пользователя
// @Tags Пользователи
// @Param updateUser body models.UpdateUserRequest true "Данные для обновления пользователя"
// @Consumes application/json
// @Produces application/json
// @Param id path int true "ID пользователя"
// @Success 204 No Content
// @Failure 422 Unprocessable Entity
// @Failure 404 Not Found
// @Router /api/users/{id} [put]
func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validation, userRequest := valiUpdatedAtUserRequest(r)
		if  validation != nil {
			http.Error(w, validation.Error(), http.StatusBadRequest)
			return
		}
		var user models.User
		user.Name = userRequest.Name
		user.Email = userRequest.Email
		user.Phone = userRequest.Phone
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])
		user.ID = paramId

		_, err := db.Exec("UPDATE users SET name = $1, email = $2, phone = $3 WHERE id = $4", user.Name, user.Email, user.Phone, paramId)
		if err != nil {
			log.Fatal(err)
		}
		errorResponse, user := getOneUser(db, paramId)
		if  errorResponse != nil {
			http.Error(w, errorResponse.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

// Документация для метода DeleteUser
// @Summary Удаляет пользователя по ID
// @Description Удаление пользователя по идентификатору
// @Tags Пользователи
// @Param id path int true "ID пользователя"
// @Success 204 No Content
// @Failure 404 Not Found
// @Router /api/users/{id} [delete]
func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])

		var user models.User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", paramId).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM users WHERE id = $1", paramId)
			if err != nil {
				//todo : fix error handling
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode("User deleted")
		}
	}
}


