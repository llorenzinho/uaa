package dtos

type CreateUserParams struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=64,password,no_whitespace,has_upper,has_lower,has_number,has_special"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type UpdateUserPasswordHashParams struct {
	Password        string `json:"password" validate:"required,min=8,max=64,password,no_whitespace,has_upper,has_lower,has_number,has_special"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type UpdateUserEmailParams struct {
	Email string `json:"email" validate:"required,email"`
}
