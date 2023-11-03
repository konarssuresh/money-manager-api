package router

import (
	"encoding/json"
	"fmt"
	"money-manager/model"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func (h handler) HandleAddFinancialTransaction(w http.ResponseWriter, r *http.Request) {
	var request model.FinancialTransactionPost
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`INSERT INTO financial_transactions (type_id,user_id,value,comment,transaction_date) VALUES (%d,'%s',%f,'%s','%s')`, request.Type_id, request.User_id, request.Value, request.Comment, request.TransactionDate)

	insertResult := h.DB.Exec(query)

	if insertResult.Error != nil {
		http.Error(w, insertResult.Error.Error(), http.StatusInternalServerError)
		return
	}

	res := model.CreateResponse{}

	if insertResult.RowsAffected > 0 {
		res.Message = "Transaction added successfully"
	} else {
		res.Message = "Transaction was not added"
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(res)
}

func (h handler) HandleGetTransactionTypeForCurrentMonth(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value("userId").(string)

	if userId == "" {
		http.Error(w, "bad request", http.StatusForbidden)
		return
	}

	query := fmt.Sprintf(`SELECT transaction_id,transaction_date,ty.type_id,type,name,value,comment
	FROM financial_transactions ft
	INNER JOIN transaction_types ty ON ft.type_id = ty.type_id
	WHERE DATE_TRUNC('month', ft.transaction_date) = DATE_TRUNC('month', current_date) AND 
	DATE_TRUNC('year',ft.transaction_date) = DATE_TRUNC('year',current_date) AND
	user_id='%s'
	ORDER BY transaction_date DESC`, userId)

	result := h.DB.Raw(query)
	var transactions []model.FinancialTransaction
	if result.Error != nil {
		http.Error(w, "Error Fetching Transactions", http.StatusInternalServerError)
		return
	}

	result.Scan(&transactions)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(transactions)

}

func (h handler) HandleGetAllFinancialTransactions(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value("userId").(string)

	if userId == "" {
		http.Error(w, "bad request", http.StatusForbidden)
		return
	}

	query := fmt.Sprintf(`SELECT transaction_id,transaction_date,ty.type_id,type,name,value,comment
	FROM financial_transactions ft
	INNER JOIN transaction_types ty ON ft.type_id = ty.type_id
	WHERE user_id='%s'
	ORDER BY transaction_date DESC`, userId)

	result := h.DB.Raw(query)
	var transactions []model.FinancialTransaction
	if result.Error != nil {
		http.Error(w, "Error Fetching Transactions", http.StatusInternalServerError)
		return
	}

	result.Scan(&transactions)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(transactions)
}

func (h handler) HandleGetFinancialTransactionsByMonthAndYear(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)

	if userId == "" {
		http.Error(w, "bad request", http.StatusForbidden)
		return
	}

	var request model.FinancialTransactionsGetByYearAndMonthPost

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`SELECT transaction_id,transaction_date,ty.type_id,type,name,value,comment,user_id
	FROM financial_transactions ft
	INNER JOIN transaction_types ty ON ft.type_id = ty.type_id
	WHERE user_id='%s'
	AND EXTRACT(YEAR FROM transaction_date) = %d
	AND EXTRACT(MONTH FROM transaction_date) = %d
	ORDER BY transaction_date DESC`, userId, request.Year, request.Month)

	var transactions []model.FinancialTransaction

	result := h.DB.Raw(query)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	result.Scan(&transactions)

	w.Header().Set("Content-Type", "application/json")

	if len(transactions) == 0 {
		// If the slice is empty, send an empty JSON array
		w.Write([]byte("[]"))
	} else {
		// Encode the non-empty slice
		json.NewEncoder(w).Encode(transactions)
	}

}

func (h handler) HandleUpdateTransaction(w http.ResponseWriter, r *http.Request) {

	var request model.FinancialTransactionPut
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var count int

	query := fmt.Sprintf(`SELECT COUNT(*) FROM financial_transactions WHERE transaction_id=%d`, request.Transaction_id)

	result := h.DB.Raw(query)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	result.Scan(&count)

	if count > 0 {
		updateQuery := fmt.Sprintf(`UPDATE financial_transactions SET type_id=%d,user_id='%s',value=%f,comment='%s',transaction_date='%s' WHERE transaction_id=%d`, request.Type_id, request.User_id, request.Value, request.Comment, request.TransactionDate, request.Transaction_id)
		updateResult := h.DB.Exec(updateQuery)

		if updateResult.Error != nil {
			http.Error(w, updateResult.Error.Error(), http.StatusInternalServerError)
			return
		}

		resp := model.CreateResponse{}

		if updateResult.RowsAffected > 0 {
			resp.Message = "Transaction updated successfully"
		} else {
			resp.Message = "Transaction not updated"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}

}

func (h handler) HandleDeleteTransaction(w http.ResponseWriter, r *http.Request) {
	var request model.RemoveTransactionDelete

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deleteQuery := fmt.Sprintf(`DELETE FROM financial_transactions WHERE transaction_id=%d`, request.Transaction_id)

	result := h.DB.Exec(deleteQuery)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	resp := model.CreateResponse{}

	if result.RowsAffected > 0 {
		resp.Message = "Transaction Deleted Successfully"
	} else {
		resp.Message = "Rows Not Deleted"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
