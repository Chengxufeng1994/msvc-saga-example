package dto

type VerifyTokenResponse struct {
	UserId    uint64 `json:"user_id"`
	IsExpired bool   `json:"is_expired"`
}
