package structs

type LoginAccount struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Remember bool   `json:"remember"`
}

type RecoveryEmail struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordToken struct {
	Token string `json:"token" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type ResetPassword struct {
	Token           string `json:"token" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required,password"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type VerifyToken struct {
	JWTToken string `json:"jwtToken" validate:"required"`
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type VerifyEmailToken struct {
	Token string `json:"token" validate:"required"`
	Email string `json:"email" validate:"required"`
}
