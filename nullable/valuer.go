package nullable

import "database/sql/driver"

func (n Nullable[T]) Value() (driver.Value, error) {
	if !n.IsValid {
		return nil, nil
	}
	return n.Data, nil
}
