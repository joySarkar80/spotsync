package dto

type Response struct {
	ID           uint   `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}
