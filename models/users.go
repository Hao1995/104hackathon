package models

// type UsersItem struct {
// 	ID    sql.NullInt64  `json:"id"`
// 	Name  sql.NullString `json:"name"`
// 	Email sql.NullString `json:"email"`
// }

type UsersItem struct {
	ID    *int64  `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
