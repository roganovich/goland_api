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

// Документация для метода GetTeams
// @Summary Возвращает список всех команд
// @Description Получение списка всех команд
// @Tags Команды
// @Accept application/json
// @Produces application/json
// @Success 200 {object} []models.TeamView
// @Failure 400 Bad Request
// @Failure 500 Internal Server Error
// @Router /api/teams [get]
func GetTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := database.DB.Query("SELECT id, name, description, city, uniform_color, participant_count, responsible_id, logo, media, status, created_at  FROM teams")
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		teams := []models.TeamView{}
		for rows.Next() {
			var teamView models.TeamView
			var responsible int
			var logo sql.NullString
			var media json.RawMessage

			if err := rows.Scan(
				&teamView.ID,
				&teamView.Name,
				&teamView.Description,
				&teamView.City,
				&teamView.UniformColor,
				&teamView.ParticipantCount,
				&responsible,
				&logo,
				&media,
				&teamView.Status,
				&teamView.CreatedAt,
				); err != nil {
				log.Println(err)
			}
			if (responsible != 0) {
				errorResponsible, responsibleUser := getUserViewById(responsible)
				if errorResponsible != nil {
					log.Println(errorResponsible.Error())
				}else{
					teamView.Responsible = responsibleUser
				}
			}
			if (!logo.Valid){
				var logoFile models.Media
				errorMedia, logoFile := getOneMedia(logo.String)
				if  errorMedia != nil {
					log.Println(errorMedia.Error())
				}else{
					teamView.Logo = &logoFile
				}
			}
			if (media == nil || len(media) == 0){
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
				teamView.Media = &mediaList
			}
			teams = append(teams, teamView)
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
		}

		json.NewEncoder(w).Encode(teams)
	}
}

func getOneTeam(paramId int64) (error, models.TeamView) {
	var team models.Team
	err := database.DB.QueryRow("SELECT * FROM teams WHERE id = $1", paramId).Scan(
		&team.ID,
		&team.Name,
		&team.Description,
		&team.City,
		&team.UniformColor,
		&team.ParticipantCount,
		&team.Responsible,
		&team.DisabilityCategory,
		&team.Logo,
		&team.Media,
		&team.Status,
		&team.CreatedAt,
		&team.UpdatedAt,
		&team.DeletedAt,
		)


	errorResponsible, responsibleUser := getUserViewById(team.Responsible)
	if  errorResponsible != nil {
		log.Println(errorResponsible.Error())
	}

	var teamView models.TeamView
	teamView.ID = team.ID
	teamView.Name = team.Name
	teamView.Description = team.Description
	teamView.City = team.City
	teamView.UniformColor = team.UniformColor
	teamView.ParticipantCount = team.ParticipantCount
	teamView.Responsible = responsibleUser
	teamView.DisabilityCategory = team.DisabilityCategory
	teamView.Status = team.Status
	teamView.CreatedAt = team.CreatedAt

	if (team.Logo != nil){
		var logoFile models.Media

		errorMedia, logoFile := getOneMedia(*team.Logo)
		if  errorMedia != nil {
			log.Println(errorMedia.Error())
		}else{
			teamView.Logo = &logoFile
		}
	}

	if (team.Media != nil){
		var mediaList []models.Media
		var mediaFiles []string

		err := json.Unmarshal(*team.Media, &mediaFiles)
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
		teamView.Media = &mediaList
	}

	return err, teamView
}

// Документация для метода GetTeam
// @Summary Возвращает информацию о команде по ID
// @Description Получение информации о команде по идентификатору
// @Tags Команды
// @Param id path int true "ID команды"
// @Success 200 {object} models.TeamView
// @Failure 400 Bad Request
// @Failure 404 Not Found
// @Router /api/teams/{id} [get]
func GetTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])

		errorResponse, teamView := getOneTeam(int64(paramId))
		if  errorResponse != nil {
			http.Error(w, errorResponse.Error(), http.StatusBadRequest)
			return
		}


		json.NewEncoder(w).Encode(teamView)
	}
}

func validateCreateTeamRequest(r *http.Request) (error, models.CreateTeamRequest) {
	var req models.CreateTeamRequest
	if validation := json.NewDecoder(r.Body).Decode(&req); validation != nil {
		return validation, req
	}
	validate := validator.New()
	if validation := validate.Struct(req); validation != nil {
		return validation, req
	}

	return nil, req
}

func validateUpdatedAtTeamRequest(r *http.Request) (error, models.UpdateTeamRequest) {
	var req models.UpdateTeamRequest
	if validation := json.NewDecoder(r.Body).Decode(&req); validation != nil {
		return validation, req
	}
	validate := validator.New()
	if validation := validate.Struct(req); validation != nil {
		return validation, req
	}

	return nil, req
}

// Документация для метода CreateTeam
// @Summary Создание новой команды
// @Description Создание новой команды
// @Tags Команды
// @Param createTeam body models.CreateTeamRequest true "Данные для создания новой команды"
// @Consumes application/json
// @Produces application/json
// @Success 201 {object} models.TeamView
// @Failure 422 Unprocessable Entity
// @Router /api/teams [post]
func CreateTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validation, teamRequest := validateCreateTeamRequest(r)
		if  validation != nil {
			http.Error(w, validation.Error(), http.StatusBadRequest)
			return
		}

		var team models.Team
		team.Name = teamRequest.Name
		team.Description = teamRequest.Description
		team.City = teamRequest.City
		team.UniformColor = teamRequest.UniformColor
		team.ParticipantCount = teamRequest.ParticipantCount

		err := database.DB.QueryRow("INSERT INTO teams (name, description, city, uniform_color, participant_count, responsible_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			team.Name,
			team.Description,
			team.City,
			team.UniformColor,
			team.ParticipantCount,
			AUTH.ID,
		).Scan(&team.ID)
		if err != nil {
			log.Println(err)
		}

		errTeam, teamView := getOneTeam(team.ID)
		if errTeam != nil {
			http.Error(w, errTeam.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(teamView)
	}
}

// Документация для метода UpdateTeam
// @Summary Обновление существующей команды
// @Description Обновление существующей команды
// @Tags Команды
// @Param updateTeam body models.UpdateTeamRequest true "Данные для обновления команды"
// @Consumes application/json
// @Produces application/json
// @Param id path int true "ID команды"
// @Success 204 No Content
// @Failure 422 Unprocessable Entity
// @Failure 404 Not Found
// @Router /api/teams/{id} [put]
func UpdateTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validation, teamRequest := validateUpdatedAtTeamRequest(r)
		if  validation != nil {
			http.Error(w, validation.Error(), http.StatusBadRequest)
			return
		}
		var team models.Team
		team.Name = teamRequest.Name
		team.Description = teamRequest.Description
		team.City = teamRequest.City
		team.Logo = teamRequest.Logo
		team.Media = teamRequest.Media
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])
		team.ID = int64(paramId)

		_, errUpdate := database.DB.Exec("UPDATE teams SET name = $1, description = $2, city = $3, logo = $4, media = $5 WHERE id = $6 and responsible_id = $7",
			team.Name,
			team.Description,
			team.City,
			team.Logo,
			team.Media,
			paramId,
			AUTH.ID)
		if errUpdate != nil {
			log.Println(errUpdate)
			http.Error(w, errUpdate.Error(), http.StatusBadRequest)

		}

		errorResponse, teamView := getOneTeam(int64(paramId))
		if  errorResponse != nil {
			http.Error(w, errorResponse.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(teamView)
	}
}

// Документация для метода DeleteTeam
// @Summary Удаляет команду по ID
// @Description Удаление команды по идентификатору
// @Tags Команды
// @Param id path int true "ID команды"
// @Success 204 No Content
// @Failure 404 Not Found
// @Router /api/teams/{id} [delete]
func DeleteTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])
		errorResponse, teamView := getOneTeam(int64(paramId))
		if errorResponse != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := database.DB.Exec("DELETE FROM teams WHERE id = $1 and responsible_id = $2", teamView.ID, AUTH.ID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode("Team deleted")
		}
	}
}


