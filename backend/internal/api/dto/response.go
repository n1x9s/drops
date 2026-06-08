package dto

type Envelope struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK(data any) Envelope {
	return Envelope{Success: true, Data: data}
}

func Fail(code string, message string) Envelope {
	return Envelope{Success: false, Error: &Error{Code: code, Message: message}}
}
