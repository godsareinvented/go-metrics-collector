package util

import "fmt"

func WrappedErrs(errors ...error) error {
	var err error
	var errWrapped error
	for _, err = range errors {
		if err == nil {
			continue
		}

		if errWrapped == nil {
			errWrapped = err
			continue
		}

		errWrapped = fmt.Errorf("%w: %w", err, errWrapped)
	}

	return errWrapped
}
