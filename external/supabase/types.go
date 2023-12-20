package supabase

import "time"

type CreateUserResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    int    `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type AppMetadata struct {
	Provider  string   `json:"provider"`
	Providers []string `json:"providers"`
}

type UserMetadata struct {
}

type IdentityData struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	PhoneVerified bool   `json:"phone_verified"`
	Sub           string `json:"sub"`
}

type Identities struct {
	IdentityID   string       `json:"identity_id"`
	ID           string       `json:"id"`
	UserID       string       `json:"user_id"`
	IdentityData IdentityData `json:"identity_data"`
	Provider     string       `json:"provider"`
	LastSignInAt time.Time    `json:"last_sign_in_at"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Email        string       `json:"email"`
}

type User struct {
	ID               string       `json:"id"`
	Aud              string       `json:"aud"`
	Role             string       `json:"role"`
	Email            string       `json:"email"`
	EmailConfirmedAt time.Time    `json:"email_confirmed_at"`
	Phone            string       `json:"phone"`
	LastSignInAt     time.Time    `json:"last_sign_in_at"`
	AppMetadata      AppMetadata  `json:"app_metadata"`
	UserMetadata     UserMetadata `json:"user_metadata"`
	Identities       []Identities `json:"identities"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type SignInUserResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    int    `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
