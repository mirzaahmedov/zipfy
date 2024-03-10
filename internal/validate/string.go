package validate

import (
	"errors"
	"regexp"
)

type StringValidatorFunc func(string) error
type StringValidator struct {
	validatorFuncs []StringValidatorFunc
	value          *string
}

func String(s *string, validatorFuncs ...StringValidatorFunc) StringValidator {
	return StringValidator{
		value: s,
	}
}

func (v StringValidator) Required(err string) StringValidator {
	v.validatorFuncs = append(v.validatorFuncs, func(s string) error {
		if s == "" {
			return errors.New(err)
		}

		return nil
	})
	return v
}
func (v StringValidator) MinLength(n int, err string) StringValidator {
	v.validatorFuncs = append(v.validatorFuncs, func(s string) error {
		if len(s) < n {
			return errors.New(err)
		}

		return nil
	})
	return v
}
func (v StringValidator) MaxLength(n int, err string) StringValidator {
	v.validatorFuncs = append(v.validatorFuncs, func(s string) error {
		if len(s) > n {
			return errors.New(err)
		}

		return nil
	})
	return v
}
func (v StringValidator) Email(err string) StringValidator {
	v.validatorFuncs = append(v.validatorFuncs, func(s string) error {
		r := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
		if !r.MatchString(s) {
			return errors.New(err)
		}
		return nil
	})
	return v
}
func (v StringValidator) Equal(s string, err string) StringValidator {
	v.validatorFuncs = append(v.validatorFuncs, func(s2 string) error {
		if s2 != s {
			return errors.New(err)
		}

		return nil
	})
	return v
}
func (v StringValidator) URL(err string) StringValidator {
	v.validatorFuncs = append(v.validatorFuncs, func(s string) error {
		r := regexp.MustCompile(`[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
		if !r.MatchString(s) {
			return errors.New(err)
		}
		return nil
	})
	return v
}
func (v StringValidator) MatchRegex(r *regexp.Regexp, err string) StringValidator {
	v.validatorFuncs = append(v.validatorFuncs, func(s string) error {
		if !r.MatchString(s) {
			return errors.New(err)
		}

		return nil
	})
	return v
}
func (v StringValidator) Validate() error {
	for _, f := range v.validatorFuncs {
		if err := f(*v.value); err != nil {
			return err
		}
	}
	return nil
}
