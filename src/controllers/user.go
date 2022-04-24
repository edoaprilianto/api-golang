package controller

import (
	"crypto/rand"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"api/src/helpers"
	"api/src/models"
)

var Name string
var Phone int64
var Role string

type UserController struct{}

type ResponseOutput struct {
	User  models.User
	Token string
}

type result struct {
	Users models.User
}

func (u UserController) SignupUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Name := r.FormValue("Name")
		Role := r.FormValue("Role")
		Phone := r.FormValue("Phone")

		// query := db.Table("users").Where("phone = ?", Phone).Find(&result)

		// if results := db.Model(&User{}).First(); results.RowsAffected > 1 {
		// 	error.ApiError(w, http.StatusNotFound, "User is already Exists!")
		// 	return
		// }

		password := randString(4)
		User := models.User{Name: Name, Phone: Phone, Role: Role, Password: password}

		if result := db.Create(&User); result.Error != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed To Add new User in database! \n"+result.Error.Error())
			return
		}

		payload := helpers.Payload{
			Name: User.Name,
			Role: User.Role,
			Id:   User.ID,
		}

		token, err := helpers.GenerateJwtToken(payload)
		if err != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed To Generate New JWT Token!")
			return
		}

		helpers.RespondWithJSON(w, ResponseOutput{
			Token: token,
			User:  User,
		})
	}

}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func (u UserController) LoginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		User := models.User{}

		Password := r.FormValue("Password")
		Phone := r.FormValue("Phone")

		if results := db.Where("phone = ? ", Phone).First(&User); results.Error != nil || results.RowsAffected < 1 {
			error.ApiError(w, http.StatusNotFound, "Invalid Phone!, Please Signup!")
			return
		}

		if User.Password != Password {
			error.ApiError(w, http.StatusNotFound, "Invalid Credentials!")
			return
		}

		payload := helpers.Payload{
			Name:  User.Name,
			Phone: User.Phone,
			Role:  User.Role,
			Id:    User.ID,
		}

		token, err := helpers.GenerateJwtToken(payload)
		if err != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed To Generate New JWT Token!")
			return
		}

		helpers.RespondWithJSON(w, ResponseOutput{
			Token: token,
			User:  User,
		})
	}
}

func (u UserController) ProfileUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		token := bearerToken[1]
		data, err := helpers.VerifyJwtToken(token)
		if err != nil {
			error.ApiError(w, http.StatusForbidden, err.Error())
			return
		}

		helpers.RespondWithJSON(w, data)
	}
}
