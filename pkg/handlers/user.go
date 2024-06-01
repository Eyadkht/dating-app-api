package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"muzz-dating/pkg/core"
	"muzz-dating/pkg/models"
	"muzz-dating/pkg/utils"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserResonse struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	// Only allow HTTP POST Method
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed: %s", r.Method)))
		return
	}

	var newUser models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err)))
		return
	}

	// Hash the password before saving to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusInternalServerError, fmt.Sprintf("Error hashing password: %v", err)))
		return
	}
	newUser.Password = string(hashedPassword)

	// Ideally the location latitude and longitude will be recieved from the frontend client
	// generating random values for now
	latitude, longitude := utils.GenerateRandomLatLong()
	newUser.Latitude = latitude
	newUser.Longitude = longitude

	result := core.GetDb().Create(&newUser)
	if err := result.Error; err != nil {
		// Check for GORM's duplicate key error
		if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				// Handle duplicate entry
				utils.WriteErrorResponse(w, utils.NewAppError(http.StatusConflict, fmt.Sprintf("User with this email address already exists: %v", newUser.Email)))
				return
			}
		} else {
			// Handle other potential errors
			utils.WriteErrorResponse(w, utils.NewAppError(http.StatusConflict, fmt.Sprintf("Error creating user: %v", err)))
			return
		}
	}

	// Used as a data transfer object to omit the Token field
	createdUserResponse := CreateUserResonse{
		ID:       newUser.ID,
		Email:    newUser.Email,
		Password: newUser.Password,
		Name:     newUser.Name,
		Gender:   newUser.Gender,
		Age:      newUser.Age,
	}
	utils.WriteSuccessResponse(w, http.StatusCreated, createdUserResponse)
}

type UserLoginResonse struct {
	Token string `json:"token"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	// Only allow HTTP POST Method
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed: %s", r.Method)))
		return
	}

	var loginPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginPayload); err != nil {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err)))
		return
	}

	var user models.User
	result := core.GetDb().Where("email = ?", loginPayload.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.WriteErrorResponse(w, utils.NewAppError(http.StatusUnauthorized, "Invalid credentials"))
			return
		} else {
			utils.WriteErrorResponse(w, utils.NewAppError(http.StatusInternalServerError, "Error retrieving user"))
			return
		}
	}

	// Compare the hashed password stored in the database with the hash of the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginPayload.Password)); err != nil {
		// Passwords do not match, return unauthorized. Returning a vauge error message for security considerations
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	// Generate token for the user and save in the database
	var generatedToken, err = createToken(&user)
	if err != nil {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusInternalServerError, fmt.Sprintf("Error generating user Token: %v", err)))
		return
	}

	userLoginResponse := UserLoginResonse{Token: generatedToken}
	utils.WriteSuccessResponse(w, http.StatusOK, userLoginResponse)
}

func createToken(user *models.User) (string, error) {

	generatedToken, err := generateTokenValue()
	if err != nil {
		return "", err
	}

	userFetchErr := core.GetDb().Preload("Token").First(&user).Error
	if userFetchErr != nil {
		return "", userFetchErr
	}

	// Check if the user already has a token
	if user.Token.ID != 0 {
		// Update the existing token value
		user.Token.Value = generatedToken
		err = core.GetDb().Save(&user.Token).Error
		if err != nil {
			return "", err
		}
		return generatedToken, nil
	} else {
		// Generate a new token for the user
		token := models.Token{
			Value:  generatedToken,
			UserID: user.ID,
		}

		// Attach the token to the user
		user.Token = token
		err = core.GetDb().Save(&user).Error
		if err != nil {
			return "", err
		}

		return generatedToken, nil
	}

}

func generateTokenValue() (string, error) {

	return uuid.New().String(), nil
}
