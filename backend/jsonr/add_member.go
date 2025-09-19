package jsonr

type AddMember struct {
	LastName         string `json:"lastName" form:"lastName"`                 // Фамилия
	FirstName        string `json:"firstName" form:"firstName"`               // Имя
	MiddleName       string `json:"middleName" form:"middleName"`             // Отчество
	DateBirth        string `json:"dateBirth" form:"dateBirth"`               // Дата рождения
	PhoneNumber      string `json:"phoneNumber" form:"phoneNumber"`           // Номер телефона
	Email            string `json:"email" form:"email"`                       // Электронная почта
	Photo            []byte `json:"photo" form:"photo"`                       // Фотография
	Password         string `json:"password" form:"password"`                 // Пароль
	MembershipStatus string `json:"membershipStatus" form:"membershipStatus"` // Статус членства
	Role             string `json:"role" form:"role"`
}
