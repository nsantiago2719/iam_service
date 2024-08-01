package models

type User struct {
	firstName, lastName string
	birthDate           string
	password            string
	email               string
	username            string
	role                Role
}
