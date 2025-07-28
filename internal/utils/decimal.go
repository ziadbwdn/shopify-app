package utils

import (
	"database/sql/driver"
	"fmt"
	"shopify-app/internal/exception"
	"strconv"

	pbdecimal "google.golang.org/genproto/googleapis/type/decimal"
)

// GormDecimal is a wrapper around pbdecimal.Decimal that implements
type GormDecimal struct {
	Internal pbdecimal.Decimal // Use a named field for the embedded struct
}

// Value implements the driver.Valuer interface for GormDecimal.
func (gd GormDecimal) Value() (driver.Value, error) {
	// Access the Value field of the internal pbdecimal.Decimal
	if gd.Internal.Value == "" {
		return nil, nil // Return nil for empty string to signify NULL in DB
	}
	// Return the string representation of the decimal value.
	return gd.Internal.Value, nil
}

// Scan implements the sql.Scanner interface for GormDecimal.
func (gd *GormDecimal) Scan(value interface{}) error {
	if value == nil {
		gd.Internal.Value = "0" // Set to "0" for null values to represent a zero decimal
		return nil
	}

	switch v := value.(type) {
	case []byte:
		gd.Internal.Value = string(v)
	case string:
		gd.Internal.Value = v
	case float64:
		// Convert float to string with sufficient precision
		gd.Internal.Value = strconv.FormatFloat(v, 'f', -1, 64) // -1 for dynamic precision
	case int64:
		gd.Internal.Value = strconv.FormatInt(v, 10)
	default:
		return fmt.Errorf("unsupported Scan type for GormDecimal: %T", value)
	}
	return nil
}

// StringToGormDecimal converts a string to a *utils.GormDecimal.
func StringToGormDecimal(s string) (*GormDecimal, *exception.AppError) {
	if s == "" {
		return &GormDecimal{Internal: pbdecimal.Decimal{Value: "0"}}, nil // Treat empty string as zero
	}
	// Basic validation: check if it's a valid number string.
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return nil, exception.NewValidationError(fmt.Sprintf("invalid decimal string format: '%s'", s), err.Error())
	}
	return &GormDecimal{Internal: pbdecimal.Decimal{Value: s}}, nil
}

// GormDecimalToString converts a *utils.GormDecimal to a string.
func GormDecimalToString(gd GormDecimal) string {
	return gd.Internal.Value
}

// MustNewGormDecimal creates a new GormDecimal from a string and panics on error.
func MustNewGormDecimal(s string) *GormDecimal {
	d, err := StringToGormDecimal(s)
	if err != nil {
		panic(err)
	}
	return d
}

// GormDecimalToFloat64 converts a GormDecimal to a float64.
func GormDecimalToFloat64(gd GormDecimal) (float64, error) {
	return strconv.ParseFloat(gd.Internal.Value, 64)
}

// Float64ToGormDecimal converts a float64 to a GormDecimal.
func Float64ToGormDecimal(f float64) (*GormDecimal, error) {
	return StringToGormDecimal(strconv.FormatFloat(f, 'f', 2, 64))
}

// GormDecimalEquals checks if two GormDecimals are equal.
func GormDecimalEquals(a, b GormDecimal) bool {
	return a.Internal.Value == b.Internal.Value
}

// MustGormDecimalAdd adds two GormDecimals and panics on error.
func MustGormDecimalAdd(a, b GormDecimal) *GormDecimal {
	fa, err := GormDecimalToFloat64(a)
	if err != nil {
		panic(err)
	}
	fb, err := GormDecimalToFloat64(b)
	if err != nil {
		panic(err)
	}
	res, err := Float64ToGormDecimal(fa + fb)
	if err != nil {
		panic(err)
	}
	return res
}
