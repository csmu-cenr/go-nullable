package nullable

import "math"

// MaxByNullableFloat64 scans items and returns the maximum valid value
// extracted by fn. Invalid values (Valid == false) and NaN are skipped.
// If no valid value exists, returns (Nullable[float64]{Valid:false}, -1).
func MaxByNullableFloat64[T any](items []T, fn func(T) Nullable[float64]) (Nullable[float64], int) {
	var (
		found bool
		best  float64
		idx   = -1
	)
	for i, it := range items {
		n := fn(it)
		if !n.Valid || math.IsNaN(n.Data) {
			continue
		}
		if !found || n.Data > best || (n.Data == best && idx == -1) {
			best = n.Data
			idx = i
			found = true
		}
	}
	if !found {
		return Nullable[float64]{Valid: false}, -1
	}
	return Nullable[float64]{Data: best, Valid: true}, idx
}

// MaxNullableFloat64 returns the max from a slice of Nullable[float64] directly.
// Same semantics as above for invalid and NaN entries.
func MaxNullableFloat64(vals []Nullable[float64]) (Nullable[float64], int) {
	return MaxByNullableFloat64(vals, func(n Nullable[float64]) Nullable[float64] { return n })
}
