package handlers

import (
	"goland_api/pkg/models"
	"database/sql"
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

// create team
func Preloader(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Загрузка файла
		file, fileHeader, errFile := r.FormFile("file")
		if errFile != nil {
			log.Fatal("Не удалось прочитать файл")
			http.Error(w, "Не удалось прочитать файл", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileName := getRandomName()
		dstPath := filepath.Join("./public/uploads/", fileName)

		f, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			log.Fatal("Не удалось открыть файл")
			http.Error(w, "Не удалось открыть файл", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		fileSize, err := io.Copy(f, file)
		if err != nil {
			log.Fatal("Не удалось скопировать файл")
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

		errInsert := db.QueryRow("INSERT INTO medias (name, path, ext, size) VALUES ($1, $2, $3, $4) RETURNING id", media.Name, media.Path, media.Ext, media.Size).Scan(&media.ID)
		if errInsert != nil {
			log.Fatal(errInsert)
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
