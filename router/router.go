package router

import (
	"fmt"
	"money-manager/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func NewDB(db *gorm.DB) handler {
	return handler{DB: db}
}

var Router *mux.Router = mux.NewRouter()

func InitializeRouter(h handler) {
	Router.Use(middleware.CorsMiddleware)
	Router.Use(middleware.AuthenticateMiddleWare)

	Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "this is a home page")
	})

	Router.HandleFunc("/signup", h.HandleSignup).Methods("POST", "OPTIONS")
	Router.HandleFunc("/login", h.HandleLogin).Methods("POST", "OPTIONS")

	Router.HandleFunc("/type/add", h.HandleAddTransactionType).Methods("POST", "OPTIONS")
	Router.HandleFunc("/type", h.HandleGetTransactionTypes).Methods("GET", "OPTIONS")
	Router.HandleFunc("/type/update", h.HandleUpdateTransactionType).Methods("PUT", "OPTIONS")
	Router.HandleFunc("/type/delete", h.HandleDeleteType).Methods("DELETE", "OPTIONS")

	Router.HandleFunc("/transaction/add", h.HandleAddFinancialTransaction).Methods("POST", "OPTIONS")
	Router.HandleFunc("/transaction/currentMonth", h.HandleGetTransactionTypeForCurrentMonth).Methods("GET", "OPTIONS")
	Router.HandleFunc("/transaction/all", h.HandleGetAllFinancialTransactions).Methods("GET", "OPTIONS")
	Router.HandleFunc("/transaction/get", h.HandleGetFinancialTransactionsByMonthAndYear).Methods("POST", "OPTIONS")
	Router.HandleFunc("/transaction/update", h.HandleUpdateTransaction).Methods("PUT", "OPTIONS")
	Router.HandleFunc("/transaction/delete", h.HandleDeleteTransaction).Methods("DELETE", "OPTIONS")

}
