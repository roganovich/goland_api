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

func getOneUser(paramId int) (error, models.UserView) {
	var userView models.UserView
	err := database.DB.QueryRow("SELECT id, name, email, phone, city, logo, media, status, created_at FROM users WHERE id = $1", int64(paramId)).Scan(
		&userView.ID,
		&userView.Name,
		&userView.Email,
		&userView.Phone,
		&userView.City,
		&userView.Logo,
		&userView.Media,
		&userView.Status,
		&userView.CreatedAt,
	)

	return err, userView
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

		errorResponse, userView := getOneUser(paramId)
		if  errorResponse != nil {
			http.Error(w, errorResponse.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(userView)
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
		errorValidation, userRequest := validateLoginUserRequest(r)
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

		tokenString, errorToken := getNewToken(user.Name, user.Email)
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

		tokenString, errorToken := getNewToken(user.Name, user.Email)
		if errorToken != nil {
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
		errValidate, userRequest := validateUpdatedAtUserRequest(r)
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

		var user models.User
		user.Name = userRequest.Name
		user.Email = userRequest.Email
		user.Phone = userRequest.Phone
		user.Password = getHashPassword(userRequest.Password)
		user.ID = userRequest.ID

		_, err := database.DB.Exec("UPDATE users SET name = $1, email = $2, phone = $3, password = $4 WHERE id = $5",
			user.Name,
			user.Email,
			user.Phone,
			user.Password,
			userRequest.ID)

		if err != nil {
			log.Println(err)
		}
		errorResponse, userView := getOneUser(userRequest.ID)
		if  errorResponse != nil {
			http.Error(w, errorResponse.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(userView)
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

	var checkUser models.CreateUserRequest
	err := database.DB.QueryRow("SELECT id, name, email, phone, password FROM users WHERE email = $1", email).Scan(
		&checkUser.ID,
		&checkUser.Name,
		&checkUser.Email,
		&checkUser.Phone,
		&checkUser.Password,
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

	var checkUser models.CreateUserRequest
	err := database.DB.QueryRow("SELECT id, name, email, phone, password FROM users WHERE phone = $1", phone).Scan(
		&checkUser.ID,
		&checkUser.Name,
		&checkUser.Email,
		&checkUser.Phone,
		&checkUser.Password,
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
		fmt.Println("Неверный формат запроса JSON")
		return nil, userRequest
	}
	vars := mux.Vars(r)
	paramId, _ := strconv.Atoi(vars["id"])
	userRequest.ID = paramId

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

// isUniqueEmailFactory создает функцию isUniqueEmail с захваченной переменной
func isUniqueEmailFactory(userRequest models.UpdateUserRequest) validator.Func {
	return func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		var checkUser models.CreateUserRequest
		err := database.DB.QueryRow("SELECT id, name, email, phone, password FROM users WHERE email = $1 AND id <> $2", email, userRequest.ID).Scan(
			&checkUser.ID,
			&checkUser.Name,
			&checkUser.Email,
			&checkUser.Phone,
			&checkUser.Password,
		)
		if err == sql.ErrNoRows {
			return true
		}
		return false
	}
}

// isUniquePhoneFactory создает функцию isUniqueEmail с захваченной переменной
func isUniquePhoneFactory(userRequest models.UpdateUserRequest) validator.Func {
	return func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		var checkUser models.CreateUserRequest
		err := database.DB.QueryRow("SELECT id, name, email, phone, password FROM users WHERE phone = $1 AND id <> $2", phone, userRequest.ID).Scan(
			&checkUser.ID,
			&checkUser.Name,
			&checkUser.Email,
			&checkUser.Phone,
			&checkUser.Password,
		)
		if err == sql.ErrNoRows {
			return true
		}
		return false
	}
}

func validateUpdatedAtUserRequest(r *http.Request) (error, models.UpdateUserRequest) {
	var userRequest models.UpdateUserRequest
	// Парсим JSON из тела запроса
	if errJson := json.NewDecoder(r.Body).Decode(&userRequest); errJson != nil {
		fmt.Println("Неверный формат запроса JSON")
		return nil, userRequest
	}
	vars := mux.Vars(r)
	paramId, _ := strconv.Atoi(vars["id"])
	userRequest.ID = paramId

	validate := validator.New()
	validate.RegisterValidation("email", isUniqueEmailFactory(userRequest))
	validate.RegisterValidation("phone", isUniquePhoneFactory(userRequest))

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

func validateLoginUserRequest(r *http.Request) (error, models.LoginUserRequest) {
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
func getNewToken(name string, email string) (string, error) {
	// ExpiresAt в миллисекундах от Unix epoch
	expiresAt := time.Now().Add(time.Hour).UnixMilli()
	claims := models.Claims{
		Username: name,
		StandardClaims: jwt.StandardClaims{
			Id:       	email,
			Subject:	name,
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
