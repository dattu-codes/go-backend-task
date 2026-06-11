package models

// CreateUserRequest validates incoming POST payloads for user creation.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	// validates that date is in exact YYYY-MM-DD format
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest validates incoming PUT payloads for user updates.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UserResponse represents the serialized schema returned to clients.
// Age is calculated dynamically on retrieval and insertion.
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"` // Serialized as YYYY-MM-DD
	Age  int    `json:"age"`
}
