package payment

type RequestNewParams struct {
	Merchant    string `json:"merchant"`
	Amount      uint   `json:"amount"`
	CallbackUrl string `json:"callbackUrl"`
	OrderId     string `json:"orderId"`
}

type ResponseNewParams struct {
	TrackId int64  `json:"trackId"`
	Result  int    `json:"result"`
	PayLink string `json:"payLink"`
	Message string `json:"message"`
}

type RequestVerifyParams struct {
	Merchant string `json:"merchant"`
	TrackId  int64  `json:"trackId"`
}

type ResponseVerifyParams struct {
	PaidAt      string `json:"paidAt"`
	Amount      uint   `json:"amount"`
	Result      int    `json:"result"`
	Status      int    `json:"status"`
	RefNumber   string `json:"refNumber"`
	Description string `json:"description"`
	CardNumber  string `json:"cardNumber"`
	OrderID     string `json:"orderId"`
	Message     string `json:"message"`
}

type CallbackParams struct {
	Success int    `schema:"success"`
	TrackId int64  `schema:"trackId"`
	OrderId string `schema:"orderId"`
	Status  int    `schema:"status"`
}
