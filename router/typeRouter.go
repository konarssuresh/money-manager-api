package router

import (
	"encoding/json"
	"fmt"
	"money-manager/model"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func (h handler) HandleAddTransactionType(w http.ResponseWriter, r *http.Request) {
	var addTypeRequest model.AddTypePost

	if err := json.NewDecoder(r.Body).Decode(&addTypeRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	validate := validator.New()

	if err := validate.Struct(addTypeRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`SELECT COUNT(*) FROM transaction_types WHERE type='%s' AND name='%s'`, addTypeRequest.Type, addTypeRequest.Name)

	result := h.DB.Raw(query)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	var count int
	result.Scan(&count)

	if count > 0 {
		http.Error(w, "This type already exists", http.StatusBadRequest)
		return
	}

	addQuery := fmt.Sprintf(`INSERT INTO transaction_types(type,name) VALUES ('%s','%s')`, addTypeRequest.Type, addTypeRequest.Name)

	updateResult := h.DB.Exec(addQuery)

	if updateResult.Error != nil {
		http.Error(w, updateResult.Error.Error(), http.StatusInternalServerError)
		return
	}

	resp := model.CreateResponse{}

	if updateResult.RowsAffected > 0 {
		resp.Message = "Type Added Successfully"
	} else {
		resp.Message = "Type was not added"
	}

	json.NewEncoder(w).Encode(resp)

}

func (h handler) HandleGetTransactionTypes(w http.ResponseWriter, r *http.Request) {
	var types []model.TransactionType

	query := `SELECT * FROM transaction_types`

	result := h.DB.Raw(query)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	result.Scan(&types)

	for _, val := range types {
		fmt.Println(val)
	}

	fmt.Println("types")

	json.NewEncoder(w).Encode(types)

}

func (h handler) HandleUpdateTransactionType(w http.ResponseWriter, r *http.Request) {
	var request model.TransactionType

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var typeCount int

	query := fmt.Sprintf(`SELECT COUNT(*) FROM transaction_types WHERE type_id=%d`, request.Type_id)

	result := h.DB.Raw(query)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	result.Scan(&typeCount)

	if typeCount > 0 {
		updateQuery := fmt.Sprintf(`UPDATE transaction_types SET type='%s',name='%s' WHERE type_id=%d`, request.Type, request.Name, request.Type_id)

		updateResult := h.DB.Exec(updateQuery)

		if updateResult.Error != nil {
			http.Error(w, updateResult.Error.Error(), http.StatusInternalServerError)
			return
		}

		resp := model.CreateResponse{}

		if updateResult.RowsAffected > 0 {
			resp.Message = "Type Updated Successfully"
		} else {
			resp.Message = "Type not updated"
		}

		json.NewEncoder(w).Encode(resp)
	}

}

func (h handler) HandleDeleteType(w http.ResponseWriter, r *http.Request) {
	var request model.RemoveTypeDelete

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`DELETE FROM transaction_types WHERE type_id=%d`, request.Type_id)

	result := h.DB.Exec(query)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	resp := model.CreateResponse{
		Message: "Type Deleted Successfully",
	}

	json.NewEncoder(w).Encode(resp)
}
