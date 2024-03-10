package validators_test

import (
	"testing"

	"github.com/carlosgonzalez/go-bundled/pkg/validators"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string `validate:"required"`
	Age  int    `validate:"gte=0"`
}

func TestValidator(t *testing.T) {

	testCases := []struct {
		name        string
		input       interface{}
		expectedErr bool
	}{
		{
			name: "valid",
			input: &User{
				Name: "John",
				Age:  30,
			},
			expectedErr: false,
		},
		{
			name: "invalid name",
			input: &User{
				Age: 30,
			},
			expectedErr: true,
		},
		{
			name: "invalid age",
			input: &User{
				Name: "John",
				Age:  -1,
			},
			expectedErr: true,
		},
	}

	validator := validators.NewCustomValidator()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.Validate(tc.input)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
