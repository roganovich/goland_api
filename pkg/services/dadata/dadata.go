package dadata

import (
	"goland_api/pkg/models"
	"bytes"
	"log"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	apiKey    = "7eef0df9113aaa894399970d7bb5f30a680b3b65"
	secretKey = "63c17bc8ddfcb0c0428f58aa66c2e12cf4f599f3"
	apiURL    = "https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address"
)

func Suggest(requestBody []byte) (models.AddressResponse, error) {
	var addressResponse models.AddressResponse

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