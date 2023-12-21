package nullable

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type NullString struct {
	Data     string
	Valid    bool
	Selected bool
}

// IsZero is the function used by the omitempty tag to determine if the field should be omitted.
func (n NullString) IsZero() bool {
	if n.Selected {
		return false
	} else {
		return true
	}
}

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	Data     int64
	Valid    bool
	Selected bool
}

// IsZero is the function used by the omitempty tag to determine if the field should be omitted.
func (n NullInt64) IsZero() bool {
	if n.Selected {
		return false
	} else {
		return true
	}
}

// NullBool is an alias for sql.NullBool data type
type NullBool struct {
	Data     bool
	Valid    bool
	Selected bool
}

// IsZero is the function used by the omitempty tag to determine if the field should be omitted.
func (n NullBool) IsZero() bool {
	if n.Selected {
		return false
	} else {
		return true
	}
}

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 struct {
	Data     float64
	Valid    bool
	Selected bool
}

// IsZero is the function used by the omitempty tag to determine if the field should be omitted.
func (n NullFloat64) IsZero() bool {
	if n.Selected {
		return false
	} else {
		return true
	}
}

type NullDate struct {
	Data     time.Time
	Valid    bool
	Selected bool
}

// IsZero is the function used by the omitempty tag to determine if the field should be omitted.
func (n NullDate) IsZero() bool {
	if n.Selected {
		return false
	} else {
		return true
	}
}

type NullTime struct {
	Data     time.Time
	Valid    bool
	Selected bool
}

// IsZero is the function used by the omitempty tag to determine if the field should be omitted.
func (n NullTime) IsZero() bool {
	if n.Selected {
		return false
	} else {
		return true
	}
}

// MarshalJSON for NullInt64
func (n NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(n.Data)
}

// UnmarshalJSON for NullInt64
func (n *NullInt64) UnmarshalJSON(b []byte) error {
	n.Selected = true
	err := json.Unmarshal(b, &n.Data)
	n.Valid = (err == nil)
	return err
}

// MarshalJSON for NullBool
func (n NullBool) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(n.Data)
}

// UnmarshalJSON for NullBool
func (n *NullBool) UnmarshalJSON(b []byte) error {
	n.Selected = true
	err := json.Unmarshal(b, &n.Data)
	n.Valid = (err == nil)
	return err
}

// MarshalJSON for NullFloat64
func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(n.Data)
}

// UnmarshalJSON for NullFloat64
func (n *NullFloat64) UnmarshalJSON(b []byte) error {
	n.Selected = true
	err := json.Unmarshal(b, &n.Data)
	n.Valid = (err == nil)
	return err
}

// MarshalJSON for NullString
func (n NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(n.Data)
}

// UnmarshalJSON for NullString
func (n *NullString) UnmarshalJSON(b []byte) error {
	n.Selected = true
	err := json.Unmarshal(b, &n.Data)
	n.Valid = (err == nil)
	return err
}

// MarshalJSON for NullTime
func (n NullTime) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	val := fmt.Sprintf("\"%s\"", n.Data.Format(time.RFC3339))
	return []byte(val), nil
}

// UnmarshalJSON for NullTime
func (n *NullTime) UnmarshalJSON(b []byte) error {
	n.Selected = true
	s := string(b)
	s = strings.ReplaceAll(s, "\"", "")

	x, err := time.Parse(time.RFC3339, s)
	if err != nil {
		n.Valid = false
		return err
	}

	n.Data = x
	n.Valid = true
	return nil
}

// MarshalJSON for NullDate
func (n NullDate) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	val := fmt.Sprintf("\"%s\"", n.Data.Format(time.DateOnly))
	return []byte(val), nil
}

// UnmarshalJSON for NullDate
func (n *NullDate) UnmarshalJSON(b []byte) error {
	n.Selected = true
	s := string(b)
	s = strings.ReplaceAll(s, "\"", "")

	x, err := time.Parse(time.DateOnly, s)
	if err != nil {
		n.Valid = false
		return err
	}

	n.Data = x
	n.Valid = true
	return nil
}
