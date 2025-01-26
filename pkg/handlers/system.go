package handlers

import (
	"goland_api/pkg/models"
	"net/http"
	"fmt"
	"os"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// AUTH - глобальная переменная с авторизацией
var AUTH *models.UserView

// Middleware для проверки авторизации
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenString := authHeader[len("Bearer "):]
		token, errToken := ParseToken(tokenString)
		if errToken != nil {
			http.Error(w, "Неверный токен", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		errorResponse, userView := getUserFromToken(token)
		if  errorResponse != nil {
			http.Error(w, "Неверный токен", http.StatusBadRequest)
			return
		}

		AUTH = userView

		// Если токен валиден, передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}

func JsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func varDump(myVar ...interface{}) {
	fmt.Printf("%v\n", myVar)
}

func dd(myVar ...interface{}) {
	varDump(myVar...)
	os.Exit(1)
}

var (
	validate *validator.Validate
	trans    ut.Translator
)

func registerCustomErrorMessages() {
	// Кастомные сообщения для стандартных тегов
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} является обязательным полем", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	validate.RegisterTranslation("min", trans, func(ut ut.Translator) error {
		return ut.Add("min", "{0} должен быть не менее {1} символов", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.Field(), fe.Param())
		return t
	})

	validate.RegisterTranslation("max", trans, func(ut ut.Translator) error {
		return ut.Add("max", "{0} должен быть не более {1} символов", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", fe.Field(), fe.Param())
		return t
	})

	validate.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} должен быть корректным email-адресом", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	// Кастомные сообщения для пользовательских тегов
	validate.RegisterTranslation("phone", trans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0} должен быть корректным номером телефона", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})
}
