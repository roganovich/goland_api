package handlers

import (
	"goland_api/pkg/models"
	"goland_api/pkg/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"fmt"
	"time"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/golang-jwt/jwt"
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
func GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := database.DB.Query("SELECT id, name, email, phone, city, logo, media, status, created_at FROM users")
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()

		users := []models.UserView{}
		for rows.Next() {
			var user models.UserView
			if err := rows.Scan(
				&user.ID,
				&user.Name,
				&user.Email,
				&user.Phone,
				&user.City,
				&user.Logo,
				&user.Media,
				&user.Status,
				&user.CreatedAt); err != nil {
				log.Println(err)
			}
			users = append(users, user)
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
		}

		json.NewEncoder(w).Encode(users)
	}
}

func getOneUser(paramId int) (error, models.User) {
	var user models.User
	err := database.DB.QueryRow("SELECT * FROM users WHERE id = $1", int64(paramId)).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.City,
		&user.Logo,
		&user.Media,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	return err, user
}

func getOneUserByEmail(paramEmail string) (error, models.CreateUserRequest) {
	var user models.CreateUserRequest
	err := database.DB.QueryRow("SELECT id, name, email, phone, password FROM users WHERE email = $1", paramEmail).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
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
func GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])

		errorResponse, user := getOneUser(paramId)
		if  errorResponse != nil {
			http.Error(w, errorResponse.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

// Документация для метода GetUser
// @Summary Возвращает информацию о пользователе по ID
// @Description Получение информации о пользователе по идентификатору
// @Tags Пользователи
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.User
// @Failure 400 Bad Request
// @Failure 404 Not Found
// @Router /api/auth/info [get]
func Info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenString := authHeader[len("Bearer "):]
		token, err := ParseToken(tokenString)
		if err != nil {
		}

		json.NewEncoder(w).Encode(token)
	}
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorValidation, userRequest := valiLoginUserRequest(r)
		if  errorValidation != nil {
			http.Error(w, errorValidation.Error(), http.StatusBadRequest)
			return
		}

		errorQuery, user := getOneUserByEmail(userRequest.Email)
		if  errorQuery != nil {
			http.Error(w, errorQuery.Error(), http.StatusBadRequest)
			return
		}
		passwordHash :=  getHashPassword(userRequest.Password)
		checkPassword := checkPasswordHash(user.Password, passwordHash)

		if checkPassword != true {
			http.Error(w, "Invalid password", http.StatusBadRequest)
			return
		}

		tokenString, errorToken := getNewToken(user)
		if errorToken != nil {
			http.Error(w, errorToken.Error(), http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(tokenString)
	}
}

func Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)
		//paramId, _ := strconv.Atoi(vars["id"])
		authHeader := r.Header.Get("Authorization")
		tokenString := authHeader[len("Bearer "):]
		token, err := ParseToken(tokenString)
		if err != nil {
			//	w.WriteHeader("WWW-Authenticate", "Bearer realm=\"Go JWT Auth\"")
			//	w.WriteHeader("Content-Type", "application/json")
			//	w.Write([]byte("Unauthorized"))
			//	return
		}

		json.NewEncoder(w).Encode(token)
		//
		//w.Write([]byte("Protected Page"))
		//
		//errorResponse, user := getOneUser(paramId)
		//if  errorResponse != nil {
		//	http.Error(w, errorResponse.Error(), http.StatusBadRequest)
		//	return
		//}
		//
		//json.NewEncoder(w).Encode(user)
	}
}

// Документация для метода CreateUser
// @Summary Создание нового пользователя
// @Description Создание нового пользователя
// @Tags Пользователи
// @Param createUser body models.CreateUserRequest true "Данные для создания пользователя"
// @Consumes application/json
// @Produces application/json
// @Success 201 {object} models.UserRegistrationResponse
// @Failure 422 Unprocessable Entity
// @Router /api/users [post]
func Registration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errValidate, userRequest := validateCreateUserRequest(r)

		if errValidate != nil {
			// Преобразуем ошибки валидации в JSON
			var validationErrors []models.ValidationErrorResponse
			for _, errValidate := range errValidate.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, models.ValidationErrorResponse{
					Field:   errValidate.Field(),
					Message: fmt.Sprintf("Ошибка в поле '%s': '%s'", errValidate.Field(), errValidate.Tag()),
				})
			}
			// Создаем структуру для ответа с ошибкой
			errorData := models.ErrorResponse{
				StatusCode:  http.StatusBadRequest,
				Message: "Возникла ошибка при регистрации",
				Errors: validationErrors,
			}
			// Сериализуем ошибки в JSON
			jsonResponse, err := json.Marshal(errorData)
			if err != nil {
				http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
				return
			}
			// Устанавливаем заголовок Content-Type
			w.Header().Set("Content-Type", "application/json")
			// Устанавливаем код состояния HTTP
			w.WriteHeader(http.StatusBadRequest)
			// Отправляем JSON-ответ
			w.Write(jsonResponse)
			return
		}

		var user models.CreateUserRequest
		user.Name = userRequest.Name
		user.Email = userRequest.Email
		user.Phone = userRequest.Phone
		user.Password = getHashPassword(userRequest.Password)

		err := database.DB.QueryRow("INSERT INTO users (name, email, phone, password) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Email, user.Phone, user.Password).Scan(&user.ID)
		if err != nil {
			http.Error(w, "Возникла ошибка при регистрации", http.StatusBadRequest)
			return
		}

		tokenString, err := getNewToken(user)
		if err != nil {
			http.Error(w, "Возникла ошибка при регистрации", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(tokenString)
		return
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
func UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validation, userRequest := validateUpdatedAtUserRequest(r)
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

		_, err := database.DB.Exec("UPDATE users SET name = $1, email = $2, phone = $3 WHERE id = $4", user.Name, user.Email, user.Phone, paramId)
		if err != nil {
			log.Println(err)
		}
		errorResponse, user := getOneUser(paramId)
		if  errorResponse != nil {
			http.Error(w, errorResponse.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(user)
		return
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
func DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])

		var user models.User
		err := database.DB.QueryRow("SELECT * FROM users WHERE id = $1", paramId).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := database.DB.Exec("DELETE FROM users WHERE id = $1", paramId)
			if err != nil {
				//todo : fix error handling
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode("User deleted")
		}
	}
}

func isUniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	//log.Println("isUniqueEmail")
	//log.Println(email)

	var user models.CreateUserRequest
	err := database.DB.QueryRow("SELECT id, name, email, phone, password FROM users WHERE email = $1", email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
	)
	if err == sql.ErrNoRows {
		return true
	}
	return false
}

func isUniquePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	//log.Println("isUniquePhone")
	//log.Println(phone)

	var user models.CreateUserRequest
	err := database.DB.QueryRow("SELECT id, name, email, phone, password FROM users WHERE phone = $1", phone).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
	)
	if err == sql.ErrNoRows {
		return true
	}
	return false
}

func validateCreateUserRequest(r *http.Request) (error, models.CreateUserRequest) {
	var userRequest models.CreateUserRequest

	// Парсим JSON из тела запроса
	if errJson := json.NewDecoder(r.Body).Decode(&userRequest); errJson != nil {
		fmt.Println("Неверный формат JSON")
		return nil, userRequest
	}

	// Вывод в формате JSON
	_, errJsonData := json.MarshalIndent(userRequest, "", "  ")
	if errJsonData != nil {
		fmt.Println("Ошибка при преобразовании в JSON:", errJsonData)
	}
	//fmt.Println(string(jsonData))
	validate := validator.New()
	validate.RegisterValidation("email", isUniqueEmail)
	validate.RegisterValidation("phone", isUniquePhone)
	errValidate := validate.Struct(userRequest)
	if errValidate != nil {
		// Если есть ошибки валидации, выводим их
		//for _, errValidate := range errValidate.(validator.ValidationErrors) {
		//	fmt.Println("Ошибка в поле '" + errValidate.Field() + "': '" + errValidate.Tag() + "'")
		//}
		return errValidate, userRequest
	} else {
		fmt.Println("Валидация прошла успешно!")
		return nil, userRequest
	}

}

func validateUpdatedAtUserRequest(r *http.Request) (error, models.UpdateUserRequest) {
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

func valiLoginUserRequest(r *http.Request) (error, models.LoginUserRequest) {
	var req models.LoginUserRequest
	if validation := json.NewDecoder(r.Body).Decode(&req); validation != nil {
		return validation, req
	}
	validate := validator.New()
	if validation := validate.Struct(req); validation != nil {
		return validation, req
	}

	return nil, req
}

// getNewToken создает новый JWT-токен
func getNewToken(user models.CreateUserRequest) (string, error) {
	// ExpiresAt в миллисекундах от Unix epoch
	expiresAt := time.Now().Add(time.Hour).UnixMilli()
	claims := models.Claims{
		Username: user.Name,
		StandardClaims: jwt.StandardClaims{
			Id:       	user.Email,
			Subject:	user.Name,
			ExpiresAt: 	expiresAt,
		},
	}
	var secretKey = []byte("my_jwt_secret")
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken проверяет и парсит JWT-токен
func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims.(*models.Claims)
		if !token.Valid || claims.Username == "" {
			return nil, fmt.Errorf("invalid token")
		}
		return claims, nil
	})

	return token, err
}

func getHashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println(err)
	}
	return string(bytes)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
