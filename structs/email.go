package structs

type VerificationTemplateData struct {
	Name             string
	VerificationLink string
}

type ResetPasswordTemplateData struct {
	Name              string
	ResetPasswordLink string
}

type NewAccountTemplateData struct {
	Name      string
	LoginLink string
	Password  string
	Email     string
}
