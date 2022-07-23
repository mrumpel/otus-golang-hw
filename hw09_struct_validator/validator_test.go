package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:10"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	AppR struct {
		Version string `validate:"regexp:\\d+"`
	}

	AppI struct {
		Version string `validate:"in:one,two,three"`
	}

	Num struct {
		Number int `validate:"min:10|max:50|in:20,60"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Password struct {
		FirstLetter  string `validate:"len:1"`
		OtherLetters string `validate:"notin:werty,23456"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in                     interface{}
		expectedErr            error
		expectedValidationErrs []error
		fails                  []string
	}{
		// no checks
		{
			in: Token{
				Header:    nil,
				Payload:   []byte("123123123123123"),
				Signature: nil,
			},
			expectedErr:            nil,
			expectedValidationErrs: nil,
		},

		// simple
		{
			in:                     App{Version: "12345"},
			expectedErr:            nil,
			expectedValidationErrs: nil,
			fails:                  nil,
		},
		{
			in:                     App{Version: "1234"},
			expectedValidationErrs: []error{errLen},
			fails:                  []string{"Version", "5"},
		},
		{
			in:                     AppR{Version: "1234567890"},
			expectedErr:            nil,
			expectedValidationErrs: nil,
			fails:                  nil,
		},
		{
			in:                     AppR{Version: "ASDF 2.5"},
			expectedErr:            nil,
			expectedValidationErrs: []error{errRegex},
			fails:                  []string{"d+", "Version"},
		},
		{
			in:                     AppI{Version: "two"},
			expectedErr:            nil,
			expectedValidationErrs: nil,
			fails:                  nil,
		},
		{
			in:                     AppI{Version: "four"},
			expectedErr:            nil,
			expectedValidationErrs: []error{errIn},
			fails:                  []string{"one,two,three", "Version"},
		},
		{
			in:                     Num{Number: 30},
			expectedErr:            nil,
			expectedValidationErrs: nil,
			fails:                  nil,
		},
		{
			in:                     Num{Number: 100},
			expectedErr:            nil,
			expectedValidationErrs: []error{errMax, errIn},
			fails:                  []string{"20,60", "50", "Number"},
		},
		{
			in:                     Num{Number: -1},
			expectedErr:            nil,
			expectedValidationErrs: []error{errMin, errIn},
			fails:                  []string{"20,60", "10", "Number"},
		},
		// complex
		{
			in: User{
				ID:     "0000000000",
				Name:   "Means Noth",
				Age:    24,
				Email:  "means@noth.io",
				Role:   "admin",
				Phones: []string{"00000000000", "11111111111"},
				meta:   nil,
			},
			expectedErr:            nil,
			expectedValidationErrs: nil,
			fails:                  nil,
		},
		{
			in: User{
				ID:     "000000000015",
				Name:   "",
				Age:    118,
				Email:  "means.noth.io",
				Role:   "nice_guy",
				Phones: []string{"000000000001", "1111"},
				meta:   nil,
			},
			expectedErr:            nil,
			expectedValidationErrs: []error{errLen, errMax, errRegex, errIn, errLen, errLen},
			fails:                  []string{"ID", "Email", "Role", "Phones"},
		},

		// test software errors
		{
			in: Response{
				Code: 250,
				Body: "<resp>",
			},
			expectedErr: errWrongValidationValue,
			fails:       []string{"404,505"},
		},
		{
			in:                     ValidationErrors{ValidationError{Err: nil, Field: "myField"}},
			expectedErr:            errNotSupportedFieldType,
			expectedValidationErrs: nil,
			fails:                  nil,
		},
		{
			in:                     Password{FirstLetter: "Q", OtherLetters: "werty"},
			expectedErr:            errWrongValidationType,
			expectedValidationErrs: nil,
			fails:                  []string{"nolen", "OtherLetters"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil && len(tt.expectedValidationErrs) == 0 {
				require.NoError(t, err)
				return
			}

			// check if software error
			if len(tt.expectedValidationErrs) == 0 {
				e := tt.expectedErr
				require.ErrorAs(t, err, &e)
				return
			}

			// check validation errors
			var l ValidationErrors
			require.ErrorAs(t, err, &l)
			require.Equal(t, len(tt.expectedValidationErrs), len(l), "Unexpected number of validation errors")

			for _, e := range tt.expectedValidationErrs {
				require.ErrorAs(t, err, &e)
			}

			// check expected words present
			for _, f := range tt.fails {
				require.ErrorContains(t, err, f)
			}

		})
	}
}
