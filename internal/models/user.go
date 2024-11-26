package models

type User struct {
	ID         int64   `json:"id"`
	Email      string  `json:"email" validate:"required,email"`
	Phone      *string `json:"phone"`
	Name       *string `json:"name"`
	LastName   *string `json:"lastName" db:"last_name"`
	SecondName *string `json:"secondName" db:"second_name"`
	BirthDate  *string `json:"birthDate" db:"birth_date" validate:"omitempty,datetime=2006-01-02"`
	Sex        *string `json:"sex" validate:"omitempty,oneof=F M"`
}
