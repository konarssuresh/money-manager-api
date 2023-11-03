package model

type TransactionType struct {
	Type_id uint   `json:"typeId"`
	Type    string `json:"type"`
	Name    string `json:"name"`
}

type FinancialTransaction struct {
	Transaction_id   uint    `json:"transactionId"`
	Transaction_date string  `json:"transactionDate"`
	Type_id          uint    `json:"typeId"`
	Type             string  `json:"type"`
	Name             string  `json:"name"`
	Value            float32 `json:"value"`
	Comment          string  `json:"comment"`
	UserId           string  `json:"userId"`
}
