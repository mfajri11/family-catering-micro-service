package HandlerHttp

type GeneralResponse struct {
	Status       string `json:"status"`
	Success      bool   `json:"success"`
	Data         any    `json:"data"`
	ResponseTime string `json:"response_time"`
}

type ErrorDetail struct {
	Description string `json:"description"`
	Code        string `json:"code"`
}

type ErrorResponse struct {
	Status       string `json:"status"`
	Success      bool   `json:"success"`
	ResponseTime string `json:"response_time"`
}
