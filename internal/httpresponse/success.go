package httpresponse

type Success struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Ok(message string, data any) Success {
	return Success{Success: true, Message: message, Data: data}
}
