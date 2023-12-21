package handler

type AppCookies struct {
	User         string
	Token        string
	RefreshToken string
	ExpiresAt    int64
	Path         string
}
