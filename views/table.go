package views

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davronkhamdamov/restaraunt_backend/models"
	"github.com/davronkhamdamov/restaraunt_backend/utils"
	"github.com/gorilla/mux"
	"github.com/skip2/go-qrcode"
)

func CreateTable(w http.ResponseWriter, r *http.Request) {
	table := models.Table{}
	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	if err := validate.Struct(table); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	dbResult := models.DB.Create(&table)
	if dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Table already exists", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusCreated, "Table created successfully", nil)
}

func GetTables(w http.ResponseWriter, r *http.Request) {
	table := []models.Table{}
	if dbResult := models.DB.Order("created_at DESC").Find(&table); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "OK", table)
}
func GetTable(w http.ResponseWriter, r *http.Request) {
	table := models.Table{}
	vars := mux.Vars(r)
	tableID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", tableID).First(&table); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Table not found", dbResult.Error.Error())
		return
	}
	url := fmt.Sprintf("https://m-menu-front.vercel.app/uz/%s", tableID)

	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "QR generation failed", err.Error())
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}
func GetOneTable(w http.ResponseWriter, r *http.Request) {
	table := models.Table{}
	vars := mux.Vars(r)
	tableID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", tableID).First(&table); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Table not found", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "ok", table)
}
func UpdateTable(w http.ResponseWriter, r *http.Request) {
	table := models.Table{}
	vars := mux.Vars(r)
	tableID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", tableID).First(&table); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Table not found", dbResult.Error.Error())
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(table); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if dbResult := models.DB.Save(&table); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Table updated successfully", nil)
}
func DeleteTable(w http.ResponseWriter, r *http.Request) {
	table := models.Table{}
	vars := mux.Vars(r)
	tableID := vars["id"]

	if dbResult := models.DB.Where("ID = ?", tableID).First(&table); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Table not found", dbResult.Error.Error())
		return
	}
	if dbResult := models.DB.Delete(&table); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Table deleted successfully", nil)
}
