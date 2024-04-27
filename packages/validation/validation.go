package validation

import (
	"fmt"

	"Edot/utilities"
)

type (
	Validation struct {
		Messages map[string]interface{}
		errors   []Error
	}

	Error struct {
		Label   string
		Message string
	}
)

// IsIntegerMin :
func (r *Validation) NewError(label string, err error) {
	r.errors = append(r.errors, Error{
		Label:   label,
		Message: err.Error(),
	})
	r.errGenerate()
}

// IsEmptyString :
func (r *Validation) IsEmptyString(value, label string) {
	if value == "" {
		err := fmt.Errorf(utilities.EmptyValue, label)
		r.errors = append(r.errors, Error{
			Label:   label,
			Message: err.Error(),
		})
		r.errGenerate()
	}
}

// IsEmailValid :
func (r *Validation) IsEmailValid(value, label string) {
	if !utilities.EmailRegex.MatchString(value) {
		err := fmt.Errorf(utilities.ValueNotValid, label)
		r.errors = append(r.errors, Error{
			Label:   label,
			Message: err.Error(),
		})
		r.errGenerate()
	}
}

// IsIntegerMax :
func (r *Validation) IsIntegerMax(value, max int, label string) {
	if value > max {
		err := fmt.Errorf(utilities.MaxValueMust, label, max)
		r.errors = append(r.errors, Error{
			Label:   label,
			Message: err.Error(),
		})
		r.errGenerate()
	}
}

// IsIntegerMin :
func (r *Validation) IsIntegerMin(value, min int, label string) {
	if value < min {
		err := fmt.Errorf(utilities.MinValueMust, label, min)
		r.errors = append(r.errors, Error{
			Label:   label,
			Message: err.Error(),
		})
		r.errGenerate()
	}
}

// IsFloatMin :
func (r *Validation) IsFloatMin(value, min float64, label string) {
	if value < min {
		err := fmt.Errorf(utilities.MinValueMust, label, min)
		r.errors = append(r.errors, Error{
			Label:   label,
			Message: err.Error(),
		})
		r.errGenerate()
	}
}

func (r *Validation) errGenerate() {
	message := map[string]interface{}{}
	for _, v := range r.errors {
		message[v.Label] = v.Message
	}
	r.Messages = message
}
