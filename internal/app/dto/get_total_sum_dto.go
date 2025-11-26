package dto

type GetTotalSum struct {
	UserID      string  `json:"user_id"`
	ServiceName *string `json:"service_name"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}
