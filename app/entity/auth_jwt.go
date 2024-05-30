package entity

// AuthJwt 認証情報を表す構造体。
type AuthJwt struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	SessionID    string `json:"session_id"`
	ExpiresIn    int    `json:"expires_in"`
}
