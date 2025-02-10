package dadata

import (
	"goland_api/pkg/models"
	"bytes"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Suggest(requestBody []byte) (models.AddressResponse, error) {
	var addressResponse models.AddressResponse
	apiKey := os.Getenv("DADATA_API_KEY")
	apiURL := os.Getenv("DADATA_API_URL")

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return addressResponse, err
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token "+apiKey)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		return addressResponse, err
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка при чтении ответа:", err)
		return addressResponse, err

	}

	// Парсим ответ
	err = json.Unmarshal(body, &addressResponse)
	if err != nil {
		log.Println("Ошибка при парсинге ответа:", err)
		return addressResponse, err
	}

	var suggestions []models.AddressSuggestion

	// Выводим результаты
	for _, suggestion := range addressResponse.Suggestions {
		//log.Println(suggestion.Value)
		suggestions = append(suggestions, suggestion)
	}
	addressResponse.Suggestions = suggestions

	return addressResponse, nil
}