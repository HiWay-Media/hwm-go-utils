package utils

type ErrorMess struct {
	Severity string `json:"severity,omitempty"`
	Message  string `json:"message,omitempty"`
}

type ApiError struct {
	Message  string `json:"message,omitempty" `
	Response string `json:"response,omitempty" `
}

type ApiMessage struct {
	Message  string      `json:"message,omitempty" `
	Response string      `json:"response,omitempty" `
	Data     interface{} `json:"data,omitempty"`
}

func ApiDefaultError(err string) *ApiError {
	return &ApiError{Response: "KO", Message: err}
}

func ApiDefaultMsgResponse(data interface{}, message string) *ApiMessage {
	return &ApiMessage{Response: "OK", Message: message, Data: data}
}

func ApiDefaultResponse(data interface{}) *ApiMessage {
	return &ApiMessage{Response: "OK", Data: data}
}

func ApiDefaultMsgOnly(message string) *ApiMessage {
	return &ApiMessage{Response: "OK", Message: message}
}
