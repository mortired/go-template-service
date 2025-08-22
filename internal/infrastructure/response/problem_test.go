package response

import (
	"net/http"
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewProblem(t *testing.T) {
	problem := NewProblem(400, "Test Title")

	assert.Equal(t, 400, problem.Status)
	assert.Equal(t, "Test Title", problem.Title)
	assert.NotNil(t, problem.Timestamp)
	assert.True(t, time.Since(problem.Timestamp) < time.Second)
}

func TestProblem_WithDetail(t *testing.T) {
	problem := NewProblem(400, "Test Title").WithDetail("Test Detail")
	assert.Equal(t, "Test Detail", problem.Detail)
}

func TestProblem_WithType(t *testing.T) {
	problem := NewProblem(400, "Test Title").WithType("test-type")
	assert.Equal(t, "test-type", problem.Type)
}

func TestProblem_WithInstance(t *testing.T) {
	problem := NewProblem(400, "Test Title").WithInstance("/test/path")
	assert.Equal(t, "/test/path", problem.Instance)
}

func TestProblem_WithTraceID(t *testing.T) {
	problem := NewProblem(400, "Test Title").WithTraceID("test-trace-id")
	assert.Equal(t, "test-trace-id", problem.TraceID)
}

func TestProblem_WithValidationErrors(t *testing.T) {
	errors := map[string][]string{
		"field1": {"error1", "error2"},
		"field2": {"error3"},
	}

	problem := NewProblem(400, "Test Title").WithValidationErrors(errors)
	assert.Equal(t, errors, problem.Errors)
}

func TestValidationError(t *testing.T) {
	validationErr := validation.Errors{
		"name":  validation.NewError("name", "Name is required"),
		"email": validation.NewError("email", "Invalid email format"),
	}

	problem := ValidationError(validationErr, "/test/path")

	assert.Equal(t, http.StatusBadRequest, problem.Status)
	assert.Equal(t, TitleValidationError, problem.Title)
	assert.Equal(t, TypeValidationError, problem.Type)
	assert.Equal(t, "/test/path", problem.Instance)
	assert.Equal(t, "Data validation failed", problem.Detail)
	assert.Equal(t, map[string][]string{
		"name":  {"Name is required"},
		"email": {"Invalid email format"},
	}, problem.Errors)
}

func TestInvalidRequest(t *testing.T) {
	problem := InvalidRequest("Test detail", "/test/path")

	assert.Equal(t, http.StatusBadRequest, problem.Status)
	assert.Equal(t, TitleInvalidRequest, problem.Title)
	assert.Equal(t, TypeInvalidRequest, problem.Type)
	assert.Equal(t, "/test/path", problem.Instance)
	assert.Equal(t, "Test detail", problem.Detail)
}

func TestInternalError(t *testing.T) {
	problem := InternalError("Test detail", "/test/path")

	assert.Equal(t, http.StatusInternalServerError, problem.Status)
	assert.Equal(t, TitleInternalError, problem.Title)
	assert.Equal(t, TypeInternalError, problem.Type)
	assert.Equal(t, "/test/path", problem.Instance)
	assert.Equal(t, "Test detail", problem.Detail)
}

func TestNotFound(t *testing.T) {
	problem := NotFound("user", "/test/path")

	assert.Equal(t, http.StatusNotFound, problem.Status)
	assert.Equal(t, "User not found", problem.Title)
	assert.Equal(t, TypeNotFound, problem.Type)
	assert.Equal(t, "/test/path", problem.Instance)
	assert.Equal(t, "Requested resource does not exist", problem.Detail)
}

func TestUnauthorized(t *testing.T) {
	problem := Unauthorized("Custom detail", "/test/path")

	assert.Equal(t, http.StatusUnauthorized, problem.Status)
	assert.Equal(t, TitleUnauthorized, problem.Title)
	assert.Equal(t, TypeUnauthorized, problem.Type)
	assert.Equal(t, "/test/path", problem.Instance)
	assert.Equal(t, "Custom detail", problem.Detail)
}

func TestForbidden(t *testing.T) {
	problem := Forbidden("Custom detail", "/test/path")

	assert.Equal(t, http.StatusForbidden, problem.Status)
	assert.Equal(t, TitleForbidden, problem.Title)
	assert.Equal(t, TypeForbidden, problem.Type)
	assert.Equal(t, "/test/path", problem.Instance)
	assert.Equal(t, "Custom detail", problem.Detail)
}

func TestConflict(t *testing.T) {
	problem := Conflict("Custom detail", "/test/path")

	assert.Equal(t, http.StatusConflict, problem.Status)
	assert.Equal(t, TitleConflict, problem.Title)
	assert.Equal(t, TypeConflict, problem.Type)
	assert.Equal(t, "/test/path", problem.Instance)
	assert.Equal(t, "Custom detail", problem.Detail)
}
