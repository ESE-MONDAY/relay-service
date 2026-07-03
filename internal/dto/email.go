package dto

type CreateEmailRequest struct {
	From    string `json:"from" validate:"required,email"`
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required,max=255"`
	Text    string `json:"text"`
	HTML    string `json:"html"`
}

type EmailResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
