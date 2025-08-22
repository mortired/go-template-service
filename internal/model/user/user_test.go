package user

import (
	"testing"
)

func TestFilter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		filter  Filter
		wantErr bool
	}{
		{
			name:    "valid filter with name",
			filter:  Filter{Name: "John"},
			wantErr: false,
		},
		{
			name:    "valid filter with ID",
			filter:  Filter{ID: func() *ID { id := ID(1); return &id }()},
			wantErr: false,
		},
		{
			name:    "valid filter with both",
			filter:  Filter{Name: "John", ID: func() *ID { id := ID(1); return &id }()},
			wantErr: false,
		},
		{
			name:    "invalid filter with zero ID",
			filter:  Filter{ID: func() *ID { id := ID(0); return &id }()},
			wantErr: true,
		},
		{
			name:    "invalid filter with negative ID",
			filter:  Filter{ID: func() *ID { id := ID(-1); return &id }()},
			wantErr: true,
		},
		{
			name:    "invalid filter with long name",
			filter:  Filter{Name: string(make([]byte, 256))},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filter.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateUserRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     CreateUserRequest{Name: "John", Email: "john@example.com"},
			wantErr: false,
		},
		{
			name:    "missing name",
			req:     CreateUserRequest{Email: "john@example.com"},
			wantErr: true,
		},
		{
			name:    "missing email",
			req:     CreateUserRequest{Name: "John"},
			wantErr: true,
		},
		{
			name:    "invalid email format",
			req:     CreateUserRequest{Name: "John", Email: "invalid-email"},
			wantErr: true,
		},
		{
			name:    "empty name",
			req:     CreateUserRequest{Name: "", Email: "john@example.com"},
			wantErr: true,
		},
		{
			name:    "name too long",
			req:     CreateUserRequest{Name: string(make([]byte, 256)), Email: "john@example.com"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUserRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
