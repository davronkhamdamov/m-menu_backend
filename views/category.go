package views

import (
	"encoding/json"
	"net/http"

	"github.com/davronkhamdamov/restaraunt_backend/models"
	"github.com/davronkhamdamov/restaraunt_backend/utils"
	"github.com/gorilla/mux"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	if err := validate.Struct(category); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if dbResult := models.DB.Create(&category); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create category", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusCreated, "Category created successfully", nil)
}
func GetCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	lang := r.URL.Query().Get("id")
	vars := mux.Vars(r)
	foodID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", foodID).First(&category); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Category not found", dbResult.Error.Error())
		return
	}
	switch lang {
	case "uz":
		category.Name = category.NameUz
	case "ru":
		category.Name = category.NameRu
	case "en":
		category.Name = category.NameEn
	default:
		category.Name = category.NameUz
	}

	utils.RespondWithSuccess(w, http.StatusOK, "OK", category)
}
func GetAllCategory(w http.ResponseWriter, r *http.Request) {
	categories := []models.Category{}
	lang := r.URL.Query().Get("lang")
	if dbResult := models.DB.Find(&categories); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch categories", dbResult.Error.Error())
		return
	}
	for i := range categories {
		switch lang {
		case "uz":
			categories[i].Name = categories[i].NameUz
		case "ru":
			categories[i].Name = categories[i].NameRu
		case "en":
			categories[i].Name = categories[i].NameEn
		default:
			categories[i].Name = categories[i].NameUz
		}
	}
	utils.RespondWithSuccess(w, http.StatusOK, "OK", categories)
}
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	vars := mux.Vars(r)
	categoryID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", categoryID).First(&category); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Category not found", dbResult.Error.Error())
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if dbResult := models.DB.Save(&category); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Category updated successfully", nil)
}
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	vars := mux.Vars(r)
	categoryID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", categoryID).First(&category); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Category not found", dbResult.Error.Error())
		return
	}
	if dbResult := models.DB.Delete(&category); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", dbResult.Error.Error())
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, "Category deleted successfully", nil)
}
