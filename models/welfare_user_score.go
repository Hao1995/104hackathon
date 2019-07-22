package models

type WelfareUserScoreDBItem struct {
	ID        *int64 `json:"id"`
	UserID    *int64 `json:"user_id"`
	WelfareNo *int64 `json:"welfare_no"`
	Score     *int64 `json:"score"`
}

type WelfareUserScoreItem struct {
	Name  *string `json:"name"`
	Score *int64  `json:"score"`
}
