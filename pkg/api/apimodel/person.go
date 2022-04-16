package apimodel

import "time"

// Person DTO is returned from a GET operation on a person.
type Person struct {
	Id             string    `json:"id" validate:"required"`             // Uniquely identifies a person
	Name           string    `json:"name" validate:"required"`           // First name, e.g. John
	Surname        string    `json:"surname" validate:"required"`        // Family name, e.g. Smith
	Email          string    `json:"email" validate:"required,email"`    // Unique email address, it can't be updated after creation
	DateOfBirth    time.Time `json:"dateOfBirth" validate:"required"`    // Date of birth formatted time.RFC3339
	Hobby          string    `json:"hobby"`                              // Favourite free-time activity
	CreatedAt      time.Time `json:"createdAt" validate:"required"`      // When this person entity was created, formatted time.RFC3339
	LastModifiedAt time.Time `json:"lastModifiedAt" validate:"required"` // The last time when this person entity was modified, formatted time.RFC3339
}

// PersonCreate DTO defines a person to create.
type PersonCreate struct {
	Name        string    `json:"name" validate:"required"`        // First name, e.g. John
	Surname     string    `json:"surname" validate:"required"`     // Family name, e.g. Smith
	Email       string    `json:"email" validate:"required,email"` // Unique email address, it can't be updated after creation
	DateOfBirth time.Time `json:"dateOfBirth" validate:"required"` // Date of birth formatted time.RFC3339, time is not relevant and may be discarded by the server
	Hobby       string    `json:"hobby" validate:"required"`       // Favourite free-time activity
}

// PersonCreated is returned after creating a person.
type PersonCreated struct {
	Id string `json:"id" validate:"required"` // Uniquely identifies a person
}

// PersonUpdate DTO defines a person to update.
type PersonUpdate struct {
	Name        string    `json:"name" validate:"required"`        // First name, e.g. John
	Surname     string    `json:"surname" validate:"required"`     // Family name, e.g. Smith
	DateOfBirth time.Time `json:"dateOfBirth" validate:"required"` // Date of birth formatted time.RFC3339, time is not relevant and may be discarded by the server
	Hobby       string    `json:"hobby" validate:"required"`       // Favourite free-time activity
}
