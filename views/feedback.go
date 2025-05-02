package views

import (
	"encoding/json"
	"net/http"

	"github.com/davronkhamdamov/restaraunt_backend/models"
	"github.com/davronkhamdamov/restaraunt_backend/utils"
	"github.com/gorilla/mux"
)

func CreateFeedback(w http.ResponseWriter, r *http.Request) {
	feedback := models.Feedback{}
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	if err := validate.Struct(feedback); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	if dbResult := models.DB.Create(&feedback); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create feedback", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusCreated, "Feedback created successfully", nil)
}
func DownloadFeedbackExcel(w http.ResponseWriter, r *http.Request) {
	var feedbacks []models.Feedback
	if err := models.DB.Preload("Table").Find(&feedbacks).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch data", err.Error())
		return
	}

	exporter := utils.NewTaomExcelExporter()

	headers := []string{"Stol raqami", "Feedback", "Mamlakat", "Yulduz", "Vaqt"}
	exporter.SetHeaders(headers)

	var rows [][]interface{}
	for _, fb := range feedbacks {
		rows = append(rows, []interface{}{
			fb.Table.Number,
			fb.Feedback,
			fb.Region,
			fb.Star,
			fb.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	exporter.SetRows(rows)

	buf, err := exporter.Generate()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="feedbacks.xlsx"`)
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func GetFeedback(w http.ResponseWriter, r *http.Request) {
	feedback := models.Feedback{}
	vars := mux.Vars(r)
	foodID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", foodID).First(&feedback); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Feedback not found", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "OK", feedback)
}
func GetAllFeedback(w http.ResponseWriter, r *http.Request) {
	categories := []models.Feedback{}
	if dbResult := models.DB.Preload("Table").Find(&categories); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch categories", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "OK", categories)
}
func DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	feedback := models.Feedback{}
	vars := mux.Vars(r)
	feedbackID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", feedbackID).First(&feedback); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Feedback not found", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Feedback deleted successfully", nil)
}
