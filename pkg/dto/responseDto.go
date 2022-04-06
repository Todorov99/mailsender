package dto

type Response struct {
	ErrorMessage string      `json:"error,omitempty"`
	MailInfo     interface{} `json:"mailInfo,omitempty"`
}
