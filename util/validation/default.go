package validation

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validatorValidate *validator.Validate
	validatorOnce     sync.Once
)

func GetValidator() *validator.Validate {
	validatorOnce.Do(func() {
		validatorValidate = validator.New()
	})

	return validatorValidate
}
