package response

import (
	"time"
)

// Problem represents error structure according to RFC 7807
type Problem struct {
	Type      string              `json:"type,omitempty"`     // URI for problem type
	Title     string              `json:"title"`              // Brief description of the problem
	Status    int                 `json:"status"`             // HTTP status code
	Detail    string              `json:"detail,omitempty"`   // Detailed description
	Instance  string              `json:"instance,omitempty"` // URI of specific instance
	TraceID   string              `json:"trace_id,omitempty"` // ID for tracing
	Timestamp time.Time           `json:"timestamp"`          // Time when error occurred
	Errors    map[string][]string `json:"errors,omitempty"`   // Validation error details
}

// NewProblem creates a new problem with basic fields
func NewProblem(status int, title string) *Problem {
	return &Problem{
		Status:    status,
		Title:     title,
		Timestamp: time.Now(),
	}
}

// WithDetail adds detailed description
func (p *Problem) WithDetail(detail string) *Problem {
	p.Detail = detail
	return p
}

// WithType adds problem type
func (p *Problem) WithType(problemType string) *Problem {
	p.Type = problemType
	return p
}

// WithInstance adds instance URI
func (p *Problem) WithInstance(instance string) *Problem {
	p.Instance = instance
	return p
}

// WithTraceID adds ID for tracing
func (p *Problem) WithTraceID(traceID string) *Problem {
	p.TraceID = traceID
	return p
}

// WithValidationErrors adds validation errors
func (p *Problem) WithValidationErrors(errors map[string][]string) *Problem {
	p.Errors = errors
	return p
}

// Common problem types
const (
	TypeValidationError = "https://tools.ietf.org/html/rfc7807#section-3.1"
	TypeInvalidRequest  = "https://tools.ietf.org/html/rfc4918#section-11.2"
	TypeInternalError   = "https://tools.ietf.org/html/rfc7231#section-6.6.1"
	TypeNotFound        = "https://tools.ietf.org/html/rfc7231#section-6.5.4"
	TypeUnauthorized    = "https://tools.ietf.org/html/rfc7235#section-3.1"
	TypeForbidden       = "https://tools.ietf.org/html/rfc7231#section-6.5.3"
	TypeConflict        = "https://tools.ietf.org/html/rfc7231#section-6.5.8"
)

// Common problem titles
const (
	TitleValidationError = "Validation error"
	TitleInvalidRequest  = "Invalid request"
	TitleInternalError   = "Server error"
	TitleNotFound        = "Resource not found"
	TitleUnauthorized    = "Unauthorized"
	TitleForbidden       = "Access denied"
	TitleConflict        = "Data conflict"
)
