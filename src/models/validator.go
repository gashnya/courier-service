package models

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

const (
	Time      = "15:04"
	Timestamp = "2006-01-02T15:04:05.000Z"
	Date      = "2006-01-02"
)

func isValidCourierType(fl validator.FieldLevel) bool {
	ct, ok := fl.Field().Interface().(CourierType)
	if !ok {
		return false
	}

	switch ct {
	case Auto, Bike, Foot:
		return true
	}

	return false
}

func isValidHours(fl validator.FieldLevel) bool {
	f := fl.Field().String()

	hours := strings.Split(f, "-")
	if len(hours) != 2 {
		return false
	}

	var err error

	t1, err := time.Parse(Time, hours[0])
	if err != nil {
		return false
	}

	t2, err := time.Parse(Time, hours[1])
	if err != nil {
		return false
	}

	return t2.After(t1)
}

func isValidTimestamp(fl validator.FieldLevel) bool {
	f := fl.Field().String()

	t, err := time.Parse(Timestamp, f)
	if err != nil {
		return false
	}

	return t.Before(time.Now())
}

func isValidDate(fl validator.FieldLevel) bool {
	f := fl.Field().String()

	_, err := time.Parse(Date, f)
	if err != nil {
		return false
	}

	return true
}

func SetupValidators(e *echo.Echo) error {
	var cv = CustomValidator{validator: validator.New()}

	if err := cv.validator.RegisterValidation("isValidCourierType", isValidCourierType); err != nil {
		return errors.New("failed to register CourierType validation")
	}
	if err := cv.validator.RegisterValidation("isValidHours", isValidHours); err != nil {
		return errors.New("failed to register hours validation")
	}
	if err := cv.validator.RegisterValidation("isValidTimestamp", isValidTimestamp); err != nil {
		return errors.New("failed to register timestamp validation")
	}
	if err := cv.validator.RegisterValidation("isValidDate", isValidDate); err != nil {
		return errors.New("failed to register date validation")
	}

	e.Validator = &cv

	return nil
}
