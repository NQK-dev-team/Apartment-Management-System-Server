package structs

type LoginAccount struct {
	Email    string `json:"email" validate:"required;"`
	Password string `json:"password" validate:"required;"`
	Remember bool   `json:"remember;"`
}

type VerifyToken struct {
	JWTToken string `json:"jwtToken" validate:"required;"`
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken" validate:"required;"`
}
