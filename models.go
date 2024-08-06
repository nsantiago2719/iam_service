package main

import (
	"github.com/jmoiron/sqlx/types"
)

var schema = `
  CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

  CREATE TABLE IF NOT EXISTS roles (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    permissions JSONB,
    name text
  );

  CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name text,
    last_name text,
    email text,
    birthdate text,
    username text,
    password text
  );

  CREATE TABLE IF NOT EXISTS users_roles (
    user_id uuid references users(id),
    role_id uuid references roles(id)
  );
  `

type User struct {
	Id        string  `json:"id" db:"users.id"`
	FirstName string  `json:"firstName" db:"users.first_name"`
	LastName  string  `json:"lastName" db:"users.last_name"`
	BirthDate string  `json:"birthDate" db:"users.birthdate"`
	Password  string  `json:"-"`
	Email     string  `json:"email" db:"users.email"`
	Username  string  `json:"userName" db:"users.username"`
	Roles     []*Role `json:"roles" db:"-"`
}

type UserRole struct {
	User
	Role
}

type Role struct {
	UserId      string         `json:"-" db:"userId"`
	Name        string         `json:"name" db:"roles.name"`
	Permissions types.JSONText `json:"permissions" db:"roles.permissions"`
}

type Permissions []Permission

type Permission struct {
	Action   string `json:"action"`
	Resource string `json:"resource"`
	Access   string `json:"access"`
}

type LoginDetails struct {
	Username string
	password string
}
