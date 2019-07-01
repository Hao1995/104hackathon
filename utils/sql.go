package utils

import (
	"database/sql"
)

func NewNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func NewNullInt64(s *int64) sql.NullInt64 {
	if s == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: *s,
		Valid: true,
	}
}

func NewNullFloat64(s *float64) sql.NullFloat64 {
	if s == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: *s,
		Valid:   true,
	}
}
