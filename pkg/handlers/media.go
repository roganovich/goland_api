package handlers

import (
	"goland_api/pkg/models"
	"goland_api/pkg/database"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"io"
	"github.com/google/uuid"
	"time"
	"path/filepath"
	"strings"
)

// get team by id
func getOneMedia(fileName string) (error, models.Media) {
	var media models.Media
	err := database.DB.QueryRow("SELECT * FROM medias WHERE name = $1", fileName).Scan(
		&media.ID,

		&media.Name,
		&media.Path,
		&media.Ext,
		&media.Size,
		&media.CreatedAt,
	)

	if  err != nil {
		log.Println(err.Error())
	}

	return err, media
}

// @Summary Загрузить медиафайл
// @Description Загрузка медиафайла
// @Tags Медиафайлы
// @Param file formData file true "Загруженный файл"
// @Success 200 {object} models.Media
// @Failure 400 {object} models.ErrorResponse
// @Failure 413 {object} models.ErrorResponse
// @Failure 415 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func Preloader() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Загрузка файла
		file, fileHeader, errFile := r.FormFile("file")
		if errFile != nil {
			log.Println("Не удалось прочитать файл")
			http.Error(w, "Не удалось прочитать файл", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileName := getRandomName()
		dstPath := filepath.Join("./public/uploads/", fileName)

		f, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			log.Println("Не удалось открыть файл")
			http.Error(w, "Не удалось открыть файл", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		fileSize, err := io.Copy(f, file)
		if err != nil {
			log.Println("Не удалось скопировать файл")
			http.Error(w, "Не удалось скопировать файл", http.StatusInternalServerError)
			return
		}

		createdAt := time.Now()
		mimeType := getMIMEType(fileHeader.Filename)

		var media models.Media
		media.Name = fileName
		media.Path = dstPath
		media.Ext = mimeType
		media.Size = fileSize
		media.CreatedAt = createdAt

		errInsert := database.DB.QueryRow("INSERT INTO medias (name, path, ext, size) VALUES ($1, $2, $3, $4) RETURNING id", media.Name, media.Path, media.Ext, media.Size).Scan(&media.ID)
		if errInsert != nil {
			log.Println(errInsert)
		}

		json.NewEncoder(w).Encode(media)
	}
}

func getRandomName() string  {
	newUUID := uuid.New()

	return newUUID.String()
}

func getMIMEType(filename string) string {
	extFile := filepath.Ext(filename)
	extData := strings.Split(extFile, ".")
	ext := ""
	if len(extData) > 0 {
		ext = extData[1]
	} else {
		log.Println("Расширение файла не удалось получить:" + extFile)
	}

	return ext
}
