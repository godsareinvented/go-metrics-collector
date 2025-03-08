package utils

import "fmt"

func WrapErrs(first, second error) error {
	if nil == second {
		return first
	}
	if nil == first {
		return second
	}
	return fmt.Errorf("%w \n%w", first, second)
}
