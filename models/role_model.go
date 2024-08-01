package models

type Role struct {
	name        string
	permissions map[string]Permission
}
