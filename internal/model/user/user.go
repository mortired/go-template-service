package user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ID int

type Filter struct {
	Name string `json:"name" query:"name"`
	ID   *ID    `json:"id" query:"id"`
}

func (f Filter) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Name, validation.Length(0, 255)),
		validation.Field(&f.ID,
			validation.NilOrNotEmpty,
			validation.When(f.ID != nil, validation.Min(1).Error("ID must be greater than 0")),
		),
	)
}

type User struct {
	ID   ID
	Name string
}

type UserResponse struct {
	ID    ID     `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UsersResponse []UserResponse

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Email string `json:"email" validate:"required,email"`
}

func (r CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required.Error("User name is required"),
			validation.Length(1, 255).Error("Name must be between 1 and 255 characters"),
		),
		validation.Field(&r.Email,
			validation.Required.Error("User email is required"),
			is.Email.Error("Invalid email format"),
		),
	)
}
