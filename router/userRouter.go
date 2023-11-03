package router

import (
	"encoding/json"
	"fmt"
	"money-manager/helper"
	"money-manager/model"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func (h handler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var signupReq model.Signup

	if err := json.NewDecoder(r.Body).Decode(&signupReq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	validate := validator.New()

	if err := validate.Struct(signupReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`SELECT COUNT(*) FROM account WHERE user_id='%s'`, signupReq.UserId)

	var count int
	result := h.DB.Raw(query).Scan(&count)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	encryptedPwd, _ := helper.Encrypt(signupReq.Password)

	insertQuery := fmt.Sprintf(`INSERT INTO account(user_id,password,created_on) VALUES ('%s','%s','%s')`, signupReq.UserId, encryptedPwd, helper.GetTodaysDate())

	insertResult := h.DB.Exec(insertQuery)

	if insertResult.Error != nil {
		http.Error(w, insertResult.Error.Error(), http.StatusInternalServerError)
	}

	response := model.CreateResponse{
		Message: "",
	}

	if insertResult.RowsAffected > 0 {
		response.Message = "User added successfully"
	} else {
		response.Message = "User was not added. Please try again"
	}

	json.NewEncoder(w).Encode(response)

}

func (h handler) HandleLogin(w http.ResponseWriter, r *http.Request) {

	fmt.Println(w.Header())
	var loginRequest model.Signup

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validate := validator.New()

	if err := validate.Struct(loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`SELECT * FROM account WHERE user_id='%s'`, loginRequest.UserId)

	result := h.DB.Raw(query)
	var user model.User
	if result.Error != nil {
		http.Error(w, "Username and password combination does not match", http.StatusUnauthorized)
		return
	}

	result.Scan(&user)

	if helper.ValidateEncryption(user.Password, loginRequest.Password) {
		updateQuery := fmt.Sprintf(`UPDATE account SET last_login='%s' WHERE user_id='%s'`, helper.GetCurrentTimeStamp(), loginRequest.UserId)
		updateResult := h.DB.Exec(updateQuery)
		if updateResult.Error != nil {
			http.Error(w, updateResult.Error.Error(), http.StatusInternalServerError)
			return
		} else {
			token, refreshToken, err := helper.GetAllToken(loginRequest.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				resp := model.LoginResponse{
					Message:      "login successful",
					Token:        token,
					RefreshToken: refreshToken,
				}
				json.NewEncoder(w).Encode(resp)
			}
		}
	} else {
		http.Error(w, "Username and password combination does not match", http.StatusUnauthorized)
		return
	}
}
