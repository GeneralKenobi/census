package model

import "time"

type Person struct {
	Id             string // Primary key
	Name           string
	Surname        string
	Email          string
	DateOfBirth    time.Time
	Hobby          string
	CreatedAt      time.Time
	LastModifiedAt time.Time
}
