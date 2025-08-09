package models

type BannedUser struct {
	ID                int
	PhoneNumber       string
	FullName          string
	Description       string
	BirthDay          string
	City              string
	SchoolFormat      string
	ApplicantUsername string
}
