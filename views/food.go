package views

import (
	"encoding/json"
	"net/http"

	"github.com/davronkhamdamov/restaraunt_backend/models"
	"github.com/davronkhamdamov/restaraunt_backend/utils"
	"github.com/gorilla/mux"
)

func CreateFood(w http.ResponseWriter, r *http.Request) {
	food := models.Food{}
	if err := json.NewDecoder(r.Body).Decode(&food); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	if err := validate.Struct(food); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	food.Available = true
	if dbResult := models.DB.Create(&food); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create food", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusCreated, "Food created successfully", nil)
}
func GetFood(w http.ResponseWriter, r *http.Request) {
	food := models.Food{}
	vars := mux.Vars(r)
	foodID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", foodID).First(&food); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get food", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "OK", food)
}
func GetAllFood(w http.ResponseWriter, r *http.Request) {
	foods := []models.Food{}
	lang := r.URL.Query().Get("lang")
	if dbResult := models.DB.Order("created_at DESC").Find(&foods); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get all food", dbResult.Error.Error())
		return
	}
	for i, food := range foods {
		switch lang {
		case "uz":
			foods[i].Name = food.NameUz
		case "ru":
			foods[i].Name = food.NameRu
		case "en":
			foods[i].Name = food.NameEn
		default:
			foods[i].Name = food.NameEn
		}
	}
	utils.RespondWithSuccess(w, http.StatusOK, "OK", foods)
}

// func GetCategoriesAndFoods(w http.ResponseWriter, r *http.Request) {
// 	var categories []models.Category
// 	lang := r.URL.Query().Get("lang")
// 	if err := models.DB.Preload("Foods").Find(&categories).Error; err != nil {
// 		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching categories and foods", err.Error())
// 		return
// 	}
// 	var validCategories []models.Category
// 	for _, category := range categories {
// 		var filteredFoods []models.Food
// 		for _, food := range category.Foods {
// 			if food.Available && food.Price > 0 && food.Image != "" {
// 				filteredFoods = append(filteredFoods, food)
// 			}
// 		}

//			if len(filteredFoods) > 0 {
//				category.Foods = filteredFoods
//				validCategories = append(validCategories, category)
//			}
//		}
//		utils.RespondWithSuccess(w, http.StatusOK, "OK", validCategories)
//	}
func GetCategoriesAndFoods(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	lang := r.URL.Query().Get("lang")
	if err := models.DB.Preload("Foods").Find(&categories).Order("created_at DESC").Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching categories and foods", err.Error())
		return
	}

	var validCategories []models.Category
	for _, category := range categories {
		var categoryNameToReturn string
		switch lang {
		case "uz":
			categoryNameToReturn = category.NameUz
		case "ru":
			categoryNameToReturn = category.NameRu
		case "en":
			categoryNameToReturn = category.NameEn
		default:
			categoryNameToReturn = category.NameEn
		}

		var filteredFoods []models.Food
		for _, food := range category.Foods {
			if food.Available && food.Price > 0 && food.ImageUrl != "" {
				var foodNameToReturn string
				var foodDescriptionReturn string
				switch lang {
				case "uz":
					foodNameToReturn = food.NameUz
					foodDescriptionReturn = food.DescriptionUz
				case "ru":
					foodNameToReturn = food.NameRu
					foodDescriptionReturn = food.DescriptionRu
				case "en":
					foodNameToReturn = food.NameEn
					foodDescriptionReturn = food.DescriptionEn
				default:
					foodNameToReturn = food.NameEn
					foodDescriptionReturn = food.DescriptionEn
				}

				filteredFood := models.Food{
					ID:          food.ID,
					Name:        foodNameToReturn,
					Description: foodDescriptionReturn,
					Price:       food.Price,
					ImageUrl:    food.ImageUrl,
					Weight:      food.Weight,
					WeightType:  food.WeightType,
					Available:   food.Available,
					CategoryID:  food.CategoryID,
				}

				filteredFoods = append(filteredFoods, filteredFood)
			}
		}

		if len(filteredFoods) > 0 {
			category.Name = categoryNameToReturn
			category.Foods = filteredFoods
			validCategories = append(validCategories, category)
		}
	}

	utils.RespondWithSuccess(w, http.StatusOK, "OK", validCategories)
}
func UpdateFood(w http.ResponseWriter, r *http.Request) {
	food := models.Food{}
	vars := mux.Vars(r)
	foodID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", foodID).First(&food); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update food", dbResult.Error.Error())
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&food); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(food); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if dbResult := models.DB.Save(&food); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "OK", nil)
}
func DeleteFood(w http.ResponseWriter, r *http.Request) {
	food := models.Food{}
	vars := mux.Vars(r)
	foodID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", foodID).First(&food); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete food", dbResult.Error.Error())
		return
	}
	if dbResult := models.DB.Delete(&food); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "OK", nil)
}
