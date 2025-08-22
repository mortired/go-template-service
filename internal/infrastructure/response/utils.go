package response

import (
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ValidationError creates problem for validation errors
func ValidationError(err error, instance string) *Problem {
	problem := NewProblem(http.StatusBadRequest, TitleValidationError).
		WithType(TypeValidationError).
		WithInstance(instance).
		WithDetail("Data validation failed")

	// Parse ozzo-validation errors
	if validationErrors, ok := err.(validation.Errors); ok {
		errors := make(map[string][]string)
		for field, fieldError := range validationErrors {
			errors[field] = []string{fieldError.Error()}
		}
		problem.WithValidationErrors(errors)
	} else {
		// If it's not a validation error, add as general detail
		problem.WithDetail(err.Error())
	}

	return problem
}

// InvalidRequest creates problem for invalid requests
func InvalidRequest(detail string, instance string) *Problem {
	return NewProblem(http.StatusBadRequest, TitleInvalidRequest).
		WithType(TypeInvalidRequest).
		WithInstance(instance).
		WithDetail(detail)
}

// InternalError creates problem for internal errors
func InternalError(detail string, instance string) *Problem {
	return NewProblem(http.StatusInternalServerError, TitleInternalError).
		WithType(TypeInternalError).
		WithInstance(instance).
		WithDetail(detail)
}

// NotFound creates problem for missing resources
func NotFound(resource string, instance string) *Problem {
	title := "Resource not found"
	if resource != "" {
		title = strings.Title(resource) + " not found"
	}

	return NewProblem(http.StatusNotFound, title).
		WithType(TypeNotFound).
		WithInstance(instance).
		WithDetail("Requested resource does not exist")
}

// Unauthorized creates problem for unauthorized requests
func Unauthorized(detail string, instance string) *Problem {
	if detail == "" {
		detail = "Authentication required"
	}

	return NewProblem(http.StatusUnauthorized, TitleUnauthorized).
		WithType(TypeUnauthorized).
		WithInstance(instance).
		WithDetail(detail)
}

// Forbidden creates problem for forbidden requests
func Forbidden(detail string, instance string) *Problem {
	if detail == "" {
		detail = "Access to resource is forbidden"
	}

	return NewProblem(http.StatusForbidden, TitleForbidden).
		WithType(TypeForbidden).
		WithInstance(instance).
		WithDetail(detail)
}

// Conflict creates problem for data conflicts
func Conflict(detail string, instance string) *Problem {
	if detail == "" {
		detail = "Data conflict"
	}

	return NewProblem(http.StatusConflict, TitleConflict).
		WithType(TypeConflict).
		WithInstance(instance).
		WithDetail(detail)
}
