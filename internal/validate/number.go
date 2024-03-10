package validate

import (
	"errors"
)

type NumberValidatorFunc func(int) error
type NumberValidator struct {
	validatorFuncs []NumberValidatorFunc
	value          *int
}

func Number(i *int, validatorFuncs ...NumberValidatorFunc) NumberValidator {
	return NumberValidator{
		value: i,
	}
}
func (v NumberValidator) NotZero(err string) NumberValidator {
	v.validatorFuncs = append(v.validatorFuncs, func(i int) error {
		if i == 0 {
			return errors.New(err)
		}

		return nil
	})
	return v
}
func (v NumberValidator) Min(n int, err string) NumberValidator {
	v.validatorFuncs = append(v.validatorFuncs, func(i int) error {
		if i < n {
			return errors.New(err)
		}

		return nil
	})
	return v
}
func (v NumberValidator) Max(n int, err string) NumberValidator {
	v.validatorFuncs = append(v.validatorFuncs, func(i int) error {
		if i > n {
			return errors.New(err)
		}

		return nil
	})
	return v
}
func (v NumberValidator) Validate() error {
	for _, f := range v.validatorFuncs {
		if err := f(*v.value); err != nil {
			return err
		}
	}
	return nil
}
