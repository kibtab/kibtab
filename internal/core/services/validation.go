// Package services holds the use cases of the core.
//
// A service calls the ports only. Release v0.2.0 fills the package.
package services

import (
	"fmt"

	"github.com/kibtab/kibtab/internal/core/domain"
)

var _ = domain.ValidationError{}

// ValidatePayload checks one sync payload before it reaches the ports.
func ValidatePayload(payload domain.SyncPayload) error {
	if payload.Table == "" {
		return domain.ValidationError{
			Row:    "",
			Field:  "",
			Value:  "",
			Reason: "table name is empty",
		}
	}
	if len(payload.Delta) == 0 {
		return domain.ValidationError{
			Row:    "",
			Field:  "",
			Value:  "",
			Reason: "delta is empty",
		}
	}
	for i, delta := range payload.Delta {
		if err := validateDelta(delta, i, payload.Delta); err != nil {
			return err
		}
	}
	return nil
}

func validateDelta(delta domain.CellDelta, index int, deltas []domain.CellDelta) error {
	if delta.Row == "" {
		return domain.ValidationError{
			Row:    delta.Row,
			Field:  delta.Field,
			Value:  delta.Value,
			Reason: fmtInvalidField("row", index),
		}
	}
	if delta.Field == "" {
		return domain.ValidationError{
			Row:    delta.Row,
			Field:  delta.Field,
			Value:  delta.Value,
			Reason: fmtInvalidField("field", index),
		}
	}
	if delta.Value == "" {
		return domain.ValidationError{
			Row:    delta.Row,
			Field:  delta.Field,
			Value:  delta.Value,
			Reason: fmtInvalidField("value", index),
		}
	}
	if delta.Version <= 0 {
		return domain.ValidationError{
			Row:    delta.Row,
			Field:  delta.Field,
			Value:  delta.Value,
			Reason: fmtInvalidField("version", index),
		}
	}
	return nil
}

func fmtInvalidField(name string, index int) string {
	return fmt.Sprintf("sync payload delta %d has an invalid %s", index, name)
}
