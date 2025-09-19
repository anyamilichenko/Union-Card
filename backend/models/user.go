package models

type Role string

const (
	ADMIN = "admin"
	USER  = "user"
)

type Accounts struct {
	Id               int    `json:"account_id" gorm:"primaryKey;autoIncrement"` // Уникальный идентификатор
	LastName         string `json:"last_name"`                                  // Фамилия
	FirstName        string `json:"first_name"`                                 // Имя
	MiddleName       string `json:"middle_name,omitempty"`                      // Отчество (может быть null)
	DateBirth        string `json:"date_birth"`                                 // Дата рождения
	PhoneNumber      string `json:"phone_number"`                               // Номер телефона
	Email            string `json:"email"`                                      // Электронная почта
	Photo            []byte `json:"photo"`                                      // Фотография (BLOB)
	HashedPassword   string `json:"password"`                                   // Пароль (захешированный)
	MembershipStatus string `json:"membership_status"`                          // Статус членства ('Active' или 'Excluded')
	Role             string `json:"role"`
}
