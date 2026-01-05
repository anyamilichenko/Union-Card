package entity

type Account struct {
	ID               uint   `json:"account_id" gorm:"primaryKey;autoIncrement"`
	LastName         string `json:"last_name"`
	FirstName        string `json:"first_name"`
	MiddleName       string `json:"middle_name,omitempty"`
	DateBirth        string `json:"date_birth"`
	PhoneNumber      string `json:"phone_number"`
	Email            string `json:"email"`
	Photo            []byte `json:"photo"`
	HashedPassword   string `json:"password"`
	MembershipStatus string `json:"membership_status"`
	Role             string `json:"role"`
}
