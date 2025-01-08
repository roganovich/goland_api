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
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Status, &user.DataCreate, &user.DateUpdate, &user.DateDelete); err != nil {
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

func getOne(db *sql.DB, paramId int) (error, models.User) {
	var user models.User
	err := db.QueryRow("SELECT * FROM users WHERE id = $1", paramId).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Status, &user.DataCreate, &user.DateUpdate, &user.DateDelete)

	return err, user
}


// get user by id
func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])

		errorResponse, user := getOne(db, paramId)
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

func validateUpdateUserRequest(r *http.Request) (error, models.UpdateUserRequest) {
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

// create user
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

// update user
func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validation, userRequest := validateUpdateUserRequest(r)
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
		errorResponse, user := getOne(db, paramId)
		if  errorResponse != nil {
			http.Error(w, errorResponse.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

// delete user
func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])

		var user models.User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", paramId).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Status, &user.DataCreate, &user.DateUpdate, &user.DateDelete)
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


