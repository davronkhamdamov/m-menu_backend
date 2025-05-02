package views

import (
	"encoding/json"
	"net/http"

	"github.com/davronkhamdamov/restaraunt_backend/models"
	"github.com/davronkhamdamov/restaraunt_backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate = validator.New()

func CreateStaff(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	if err := validate.Struct(user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	var err error
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}
	if dbResult := models.DB.Create(&user); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusCreated, "User created successfully", nil)
}
func GetStaffs(w http.ResponseWriter, r *http.Request) {
	var staff []models.User
	if dbResult := models.DB.Find(&staff); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Staff retrieved successfully", staff)
}
func GetStaff(w http.ResponseWriter, r *http.Request) {
	var staff models.User
	vars := mux.Vars(r)
	tableID := vars["id"]
	if dbResult := models.DB.Where("ID = ?", tableID).First(&staff); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Staff retrieved successfully", staff)
}
func UpdateStaff(w http.ResponseWriter, r *http.Request) {
	var existingStaff models.User
	vars := mux.Vars(r)
	tableID := vars["id"]

	if dbResult := models.DB.Where("ID = ?", tableID).First(&existingStaff); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", dbResult.Error.Error())
		return
	}

	var updatedData models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(updatedData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingStaff.Login = updatedData.Login
	existingStaff.Role = updatedData.Role

	if updatedData.Password != "" {
		hashedPassword, err := utils.HashPassword(updatedData.Password)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", err.Error())
			return
		}
		existingStaff.Password = hashedPassword
	}

	if dbResult := models.DB.Save(&existingStaff); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", dbResult.Error.Error())
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, "Staff updated successfully", existingStaff)
}
func Login(w http.ResponseWriter, r *http.Request) {
	user := models.Login{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	if err := validate.Struct(user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	var dbUser models.User
	if dbResult := models.DB.Where("login = ?", user.Login).First(&dbUser); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid login or password", "")
		return
	}
	if err := utils.CheckPassword(dbUser.Password, user.Password); err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid login or password", "")
		return
	}
	token, err := utils.CreateToken(dbUser.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Login successful", map[string]string{"role": dbUser.Role.String(), "token": token})
}
