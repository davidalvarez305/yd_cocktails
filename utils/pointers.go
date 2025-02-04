package utils

import "database/sql"

func CreateNullString(ptr *string) sql.NullString {
	if ptr == nil || *ptr == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *ptr, Valid: true}
}

func CreateNullFloat64(ptr *float64) sql.NullFloat64 {
	if ptr == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: *ptr, Valid: true}
}

func CreateNullInt64(ptr *int64) sql.NullInt64 {
	if ptr == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *ptr, Valid: true}
}

func CreateNullInt(ptr *int) sql.NullInt64 {
	if ptr == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*ptr), Valid: true}
}

func CreateNullBool(ptr *bool) sql.NullBool {
	if ptr == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *ptr, Valid: true}
}

func CreateNullBoolDefaultFalse(ptr *bool) sql.NullBool {
	if ptr == nil {
		return sql.NullBool{Bool: false, Valid: true}
	}
	return sql.NullBool{Bool: *ptr, Valid: true}
}
