package helpers

type (
	ValidationError struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}
)

type Response struct {
	Status     int               `json:"status"`
	Message    string            `json:"message"`
	Validation []ValidationError `json:"validation,omitempty"`
	Data       interface{}       `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
}

func BuildResponse(data *Response) map[string]interface{} {
	return map[string]interface{}{
		"status":     data.Status,
		"message":    data.Message,
		"validation": data.Validation,
		"data":       data.Data,
	}
}
