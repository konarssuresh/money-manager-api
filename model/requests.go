package model

type Signup struct {
	UserId   string `json:"userId" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AddTypePost struct {
	Type string `json:"type" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type RemoveTypeDelete struct {
	Type_id int `json:"typeId" validate:"required"`
}

type FinancialTransactionPost struct {
	Type_id         int     `json:"typeId" validate:"required"`
	User_id         string  `json:"userId" validate:"required"`
	Value           float32 `json:"value" validate:"required,gte=1"`
	Comment         string  `json:"comment" validate:"required"`
	TransactionDate string  `json:"transactionDate" validate:"required,datetime=2006-01-02"`
}

type FinancialTransactionPut struct {
	Transaction_id  int     `json:"transactionId" validate:"required"`
	Type_id         int     `json:"typeId" validate:"required"`
	User_id         string  `json:"userId" validate:"required"`
	Value           float32 `json:"value" validate:"required,gte=1"`
	Comment         string  `json:"comment" validate:"required"`
	TransactionDate string  `json:"transactionDate" validate:"required,datetime=2006-01-02"`
}

type FinancialTransactionsGetByYearAndMonthPost struct {
	Year  int `json:"year" validate:"required,gte=1900,lte=9999"`
	Month int `json:"month" validate:"required,gte=1,lte=12"`
}

type RemoveTransactionDelete struct {
	Transaction_id int `json:"transactionId" validate:"required"`
}
