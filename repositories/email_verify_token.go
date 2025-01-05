package repositories

type EmailVerifyTokenRepository struct{}

func NewEmailVerifyTokenRepository() *EmailVerifyTokenRepository {
	return &EmailVerifyTokenRepository{}
}
