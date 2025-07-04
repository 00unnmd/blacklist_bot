package models

type BannedUser struct {
	ID           int
	PhoneNumber  string
	FullName     string
	Description  string
	BirthDay     string
	City         string
	SchoolFormat string
}

type Appeal struct {
	ID         int
	Question   string
	Initiator  string
	IsAnswered bool
}
