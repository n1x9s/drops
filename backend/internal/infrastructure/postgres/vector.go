package postgres

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Vector []float32

func (v Vector) Value() (driver.Value, error) {
	if len(v) == 0 {
		return nil, nil
	}
	parts := make([]string, len(v))
	for i, f := range v {
		parts[i] = strconv.FormatFloat(float64(f), 'f', -1, 32)
	}
	return "[" + strings.Join(parts, ",") + "]", nil
}

func (v *Vector) Scan(value interface{}) error {
	if value == nil {
		*v = nil
		return nil
	}
	var raw string
	switch typed := value.(type) {
	case string:
		raw = typed
	case []byte:
		raw = string(typed)
	default:
		return fmt.Errorf("scan vector: unsupported type %T", value)
	}
	raw = strings.Trim(raw, "[]")
	if raw == "" {
		*v = nil
		return nil
	}
	fields := strings.Split(raw, ",")
	out := make([]float32, 0, len(fields))
	for _, field := range fields {
		f, err := strconv.ParseFloat(strings.TrimSpace(field), 32)
		if err != nil {
			return err
		}
		out = append(out, float32(f))
	}
	*v = out
	return nil
}

func (Vector) GormDataType() string {
	return "vector"
}

func (Vector) GormDBDataType(_ *gorm.DB, _ *schema.Field) string {
	return "vector(768)"
}
