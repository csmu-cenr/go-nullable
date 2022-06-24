package nullable

import (
	"database/sql"
	"fmt"
	"time"
)

func (n *Nullable[T]) getScanner() sql.Scanner {
	dest := any(&n.Data)
	switch s := dest.(type) {
	case *string:
		return &sql.NullString{String: *s, Valid: n.IsValid}
	case *bool:
		return &sql.NullBool{Bool: *s, Valid: n.IsValid}
	case *float64:
		return &sql.NullFloat64{Float64: *s, Valid: n.IsValid}
	case *float32:
		return &sql.NullFloat64{Float64: float64(*s), Valid: n.IsValid}
	case *int16:
		return &sql.NullInt16{Int16: *s, Valid: n.IsValid}
	case *int32:
		return &sql.NullInt32{Int32: *s, Valid: n.IsValid}
	case *int:
		return &sql.NullInt32{Int32: int32(*s), Valid: n.IsValid}
	case *int64:
		return &sql.NullInt64{Int64: *s, Valid: n.IsValid}
	case *time.Time:
		return &sql.NullTime{Time: *s, Valid: n.IsValid}
	}

	return nil
}

func (n *Nullable[T]) getScannerValue(scanner sql.Scanner) any {
	switch s := scanner.(type) {
	case *sql.NullString:
		n.IsValid = s.Valid
		n.Data = any(s.String).(T)
	case *sql.NullBool:
		n.IsValid = s.Valid
		n.Data = any(s.Bool).(T)
	case *sql.NullFloat64:
		n.IsValid = s.Valid
		switch any(n.Data).(type) {
		case float32:
			n.Data = any(float32(s.Float64)).(T)
		case float64:
			n.Data = any(s.Float64).(T)
		}
	case *sql.NullInt16:
		n.IsValid = s.Valid
		n.Data = any(s.Int16).(T)
	case *sql.NullInt32:
		n.IsValid = s.Valid
		switch any(n.Data).(type) {
		case int:
			n.Data = any(int(s.Int32)).(T)
		case int32:
			n.Data = any(s.Int32).(T)
		}
	case *sql.NullInt64:
		n.IsValid = s.Valid
		n.Data = any(s.Int64).(T)
	case *sql.NullTime:
		n.IsValid = s.Valid
		n.Data = any(s.Time).(T)
	}

	return nil
}

func (n *Nullable[T]) Scan(value any) error {
	if value == nil {
		n.IsValid = false
		return nil
	}

	scanner := n.getScanner()
	if scanner == nil {
		n.IsValid = false
		var ref T
		return fmt.Errorf("no scanner available for %T", ref)
	}
	err := scanner.Scan(value)
	if err != nil {
		n.IsValid = false
		return err
	}

	n.getScannerValue(scanner)
	return nil
}
