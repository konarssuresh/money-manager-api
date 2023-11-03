package model

type CreateResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Message      string `json:"message"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type TypesResponse []TransactionType
