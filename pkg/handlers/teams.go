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
func GetTeams(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM teams")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		teams := []models.Team{}
		for rows.Next() {
			var team models.Team
			if err := rows.Scan(
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
				); err != nil {
				log.Fatal(err)
			}
			teams = append(teams, team)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(teams)
	}
}

func getOneTeam(db *sql.DB, paramId int) (error, models.TeamView) {
	var team models.Team
	err := db.QueryRow("SELECT * FROM teams WHERE id = $1", int64(paramId)).Scan(
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


	var teamView models.TeamView
	teamView.ID = team.ID
	teamView.Name = team.Name
	teamView.Description = team.Description
	teamView.City = team.City
	teamView.UniformColor = team.UniformColor
	teamView.ParticipantCount = team.ParticipantCount
	teamView.Responsible = team.Responsible
	teamView.DisabilityCategory = team.DisabilityCategory
	teamView.Status = team.Status
	teamView.CreatedAt = team.CreatedAt
	teamView.UpdatedAt = team.UpdatedAt
	teamView.DeletedAt = team.DeletedAt

	if (team.Logo != nil){
		var logoFile models.Media

		errorMedia, logoFile := getOneMedia(db, *team.Logo)
		if  errorMedia != nil {
			log.Fatal(errorMedia.Error())
		}else{
			teamView.Logo = &logoFile
		}
	}

	if (team.Media != nil){
		var mediaList []models.Media
		var mediaFiles []string

		err := json.Unmarshal(*team.Media, &mediaFiles)
		if err != nil {
			log.Fatal("Ошибка при парсинге JSON:", err)
		}
		for _, mediaFile := range mediaFiles {
			errorMedia, mediaFile := getOneMedia(db, mediaFile)
			if  errorMedia != nil {
				log.Fatal(errorMedia.Error())
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
func GetTeam(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])

		errorResponse, teamView := getOneTeam(db, paramId)
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

func valiUpdatedAtTeamRequest(r *http.Request) (error, models.UpdateTeamRequest) {
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
func CreateTeam(db *sql.DB) http.HandlerFunc {
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

		err := db.QueryRow("INSERT INTO teams (name, description, city) VALUES ($1, $2, $3) RETURNING id", team.Name, team.Description, team.City).Scan(&team.ID)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(team)
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
func UpdateTeam(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validation, teamRequest := valiUpdatedAtTeamRequest(r)
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

		_, err := db.Exec("UPDATE teams SET name = $1, description = $2, city = $3, logo = $4, media = $5 WHERE id = $6",
			team.Name,
			team.Description,
			team.City,
			team.Logo,
			team.Media,
			paramId)
		if err != nil {
			log.Fatal(err)
		}
		errorResponse, teamView := getOneTeam(db, paramId)
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
func DeleteTeam(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramId, _ := strconv.Atoi(vars["id"])

		var team models.Team
		err := db.QueryRow("SELECT * FROM teams WHERE id = $1", paramId).Scan(
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
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM teams WHERE id = $1", paramId)
			if err != nil {
				//todo : fix error handling
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode("Team deleted")
		}
	}
}


