package dto

type Response struct {
	ErrorMessage    string `json:"error,omitempty"`
	ResponseMessage string `json:"responseMessage,omitempty"`
}
