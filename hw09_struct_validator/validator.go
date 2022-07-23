package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	b := strings.Builder{}
	for _, e := range v {
		b.WriteString(e.Error())
		b.WriteRune('\n')
	}
	return b.String()
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%v in field %v", e.Err, e.Field)
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

const validateTagKey = "validate"

// software errors.
var (
	errNotAStruct            = errors.New("not a struct")
	errNotSupportedFieldType = errors.New("not supported field type")
	errWrongValidationType   = errors.New("wrong validation type ")
	errWrongValidationValue  = errors.New("wrong validation value ")
)

// validation errors.
var (
	errValidation = errors.New("invalid")
	errMin        = fmt.Errorf("%w min: less then", errValidation)
	errMax        = fmt.Errorf("%w max: greater then", errValidation)
	errIn         = fmt.Errorf("%w in: not in", errValidation)
	errLen        = fmt.Errorf("%w len: not equal", errValidation)
	errRegex      = fmt.Errorf("%w regexp: not match with", errValidation)
)

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		return errNotAStruct
	}

	valErrors := make(ValidationErrors, 0)
	for fnum := 0; fnum < typ.NumField(); fnum++ {
		tagstring, ok := typ.Field(fnum).Tag.Lookup(validateTagKey)
		if !ok {
			continue
		}

		tags := strings.Split(tagstring, "|")
		for _, tag := range tags {
			valtag := strings.SplitN(tag, ":", 2)

			err := validateTag(val.Field(fnum), typ.Field(fnum).Name, valtag[0], valtag[1])

			var valTagErrors ValidationErrors
			var valTagError ValidationError

			if errors.As(err, &valTagError) {
				valErrors = append(valErrors, valTagError)
				continue
			}

			if errors.As(err, &valTagErrors) {
				valErrors = append(valErrors, valTagErrors...)
				continue
			}

			if err != nil {
				return err
			}
		}
	}

	if len(valErrors) != 0 {
		return valErrors
	}
	return nil
}

func validateTag(v reflect.Value, name, valType, valValue string) error {
	var err error

	switch v.Kind() { //nolint:exhaustive
	case reflect.Int:
		f, makeErr := makeIntValidator(name, valType, valValue)
		if makeErr != nil {
			return makeErr
		}
		t := int(v.Int())
		err = f(t)

	case reflect.String:
		f, makeErr := makeStringValidator(name, valType, valValue)
		if makeErr != nil {
			return makeErr
		}
		t := v.String()

		err = f(t)

	case reflect.Slice:

		if v.Type().Elem().Kind() == reflect.Int {
			f, makeErr := makeIntValidator(name, valType, valValue)
			if makeErr != nil {
				return makeErr
			}
			err = validateIntSlice(v, f)
			break
		}

		if v.Type().Elem().Kind() == reflect.String {
			f, makeErr := makeStringValidator(name, valType, valValue)
			if makeErr != nil {
				return makeErr
			}
			err = validateStringSlice(v, f)
			break
		}

		err = errNotSupportedFieldType

	default:
		err = errNotSupportedFieldType
	}

	return err
}

func validateIntSlice(v reflect.Value, f func(int) error) error {
	var errs ValidationErrors
	for i := 0; i < v.Len(); i++ {
		x := int(v.Index(i).Int())
		err := f(x)
		if err == nil {
			continue
		}

		var e ValidationError
		ok := errors.As(err, &e)
		if !ok {
			return err
		}
		e.Err = fmt.Errorf("%w at position %v", e.Err, i)
		errs = append(errs, e)
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func validateStringSlice(v reflect.Value, f func(string) error) error {
	var errs ValidationErrors
	for i := 0; i < v.Len(); i++ {
		x := v.Index(i).String()
		err := f(x)
		if err == nil {
			continue
		}

		var e ValidationError
		ok := errors.As(err, &e)
		if !ok {
			return err
		}
		e.Err = fmt.Errorf("%w at position %v", e.Err, i)
		errs = append(errs, e)
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func makeIntValidator(name, tag, value string) (func(int) error, error) {
	var f func(int) error
	switch tag {
	case "min":
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("%w tag %v value %v", errWrongValidationValue, tag, value)
		}

		f = func(i int) error {
			if i < intvalue {
				return ValidationError{
					Field: name,
					Err:   fmt.Errorf("%w %v", errMin, value),
				}
			}
			return nil
		}

	case "max":
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("%w tag %v value %v", errWrongValidationValue, tag, value)
		}

		f = func(i int) error {
			if i > intvalue {
				return ValidationError{
					Field: name,
					Err:   fmt.Errorf("%w %v", errMax, value),
				}
			}
			return nil
		}

	case "in":
		values := strings.SplitN(value, ",", 2)
		min, err := strconv.Atoi(values[0])
		if err != nil {
			return nil, fmt.Errorf("%w tag %v value %v", errWrongValidationValue, tag, value)
		}
		max, err := strconv.Atoi(values[1])
		if err != nil {
			return nil, fmt.Errorf("%w tag %v value %v", errWrongValidationValue, tag, value)
		}

		f = func(i int) error {
			if i < min || i > max {
				return ValidationError{
					Field: name,
					Err:   fmt.Errorf("%w %v", errIn, value),
				}
			}
			return nil
		}

	default:
		return nil, fmt.Errorf("%w %v", errWrongValidationType, value)
	}

	return f, nil
}

func makeStringValidator(name, tag, value string) (func(string) error, error) {
	var f func(string) error
	switch tag {
	case "len":
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("%w tag %v value %v", errWrongValidationValue, tag, value)
		}

		f = func(s string) error {
			if len(s) != intvalue {
				return ValidationError{
					Field: name,
					Err:   fmt.Errorf("%w %v", errLen, value),
				}
			}
			return nil
		}

	case "regexp":
		re, err := regexp.Compile("^" + value + "$") // full string matching
		if err != nil {
			return nil, fmt.Errorf("%w, regexp not compiled: %v", errWrongValidationValue, value)
		}

		f = func(s string) error {
			if re.MatchString(s) {
				return nil
			}

			fmt.Println("")
			return ValidationError{
				Field: name,
				Err:   fmt.Errorf("%w %v", errRegex, value),
			}
		}

	case "in":
		values := strings.Split(value, ",")

		f = func(s string) error {
			for i := range values {
				if values[i] == s {
					return nil
				}
			}

			return ValidationError{
				Field: name,
				Err:   fmt.Errorf("%w %v", errIn, value),
			}
		}

	default:
		return nil, fmt.Errorf("%w %v", errWrongValidationType, value)
	}

	return f, nil
}
