package models

type User struct {
	FirstName, LastName string
	BirthDate           string
	Password            string
	Email               string
	Username            string
	UserRole            Role
}
