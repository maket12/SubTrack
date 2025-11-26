package dto

type GetSubscriptionList struct {
	UserID      *string `json:"user_id"`
	ServiceName *string `json:"service_name"`
	Limit       int     `json:"limit"`
	Offset      int     `json:"offset"`
}
