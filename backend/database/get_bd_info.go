package database

type UserInfo struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	MiddleName  *string `json:"middle_name,omitempty"`
	DateBirth   string  `json:"date_birth"`
	PhoneNumber string  `json:"phone_number"`
	Email       string  `json:"email"`
	Photo       []byte  `json:"photo"`
}
