package conv

import (
	"database/sql"
	"strings"
	"time"
)

//ToString конвертор типа из NullString в string
func ToString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}

//ToInt64 конвертор типа из NullInt64 в int64
func ToInt64(value sql.NullInt64) int64 {
	if value.Valid {
		return value.Int64
	}
	return 0
}

//ToInt32 конвертор типа из NullInt32 в int32
func ToInt32(value sql.NullInt32) int32 {
	if value.Valid {
		return value.Int32
	}
	return 0
}

//ToFloat64 конвертор типа из NullFloat64 в float64
func ToFloat64(value sql.NullFloat64) float64 {
	if value.Valid {
		return value.Float64
	}
	return 0
}

//ToTime конвертор типа из NullTime в Time
func ToTime(value sql.NullTime) time.Time {
	if value.Valid {
		return value.Time
	}
	return time.Time{}
}

//ToTimeString конвертор типа из NullTime в string по формату: '2006-01-02'
func ToTimeString(value sql.NullTime) string {
	if value.Valid {
		return value.Time.Format("2006-01-02")
	}
	return ""
}

//ToBool конвертор типа из NullBool в bool
func ToBool(value sql.NullBool) bool {
	if value.Valid {
		return value.Bool
	}
	return false
}

//ToSQLString конвертор типа из string в NullString
func ToSQLString(value string) sql.NullString {
	sValue := strings.TrimSpace(value)
	if sValue == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: sValue, Valid: true}
}

//ToSQLInt32 конвертор типа из int32 в NullInt32
func ToSQLInt32(value int32) sql.NullInt32 {
	if value == 0 {
		return sql.NullInt32{Int32: 0, Valid: false}
	}
	return sql.NullInt32{Int32: value, Valid: true}
}

//ToSQLInt32Def0 конвертор типа из int32 в NullInt32
//Если value = -1, то в БД пишем NULL
//Если value = 0, то в БД пишем 0
func ToSQLInt32Def0(value int32) sql.NullInt32 {
	if value == -1 {
		return sql.NullInt32{Int32: -1, Valid: false}
	}
	return sql.NullInt32{Int32: value, Valid: true}
}

//ToSQLInt64 конвертор типа из int64 в NullInt64
func ToSQLInt64(value int64) sql.NullInt64 {
	if value == 0 {
		return sql.NullInt64{Int64: 0, Valid: false}
	}
	return sql.NullInt64{Int64: value, Valid: true}
}

//ToSQLTime конвертор типа из Time в NullTime
func ToSQLTime(value time.Time) sql.NullTime {
	if value.IsZero() {
		return sql.NullTime{Time: time.Time{}, Valid: false}
	}
	return sql.NullTime{Time: value, Valid: true}
}

//ToSQLBool конвертор типа из bool в NullBool
func ToSQLBool(value bool) sql.NullBool {
	if !value {
		return sql.NullBool{Bool: false, Valid: false}
	}
	return sql.NullBool{Bool: value, Valid: true}
}