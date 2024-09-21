package main

import (
	"github.com/golang-jwt/jwt/v5"
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
    email text,
    username text,
    password text
  );

  CREATE TABLE IF NOT EXISTS users_roles (
    user_id uuid references users(id),
    role_id uuid references roles(id)
  );
  `

// User used for creating struct for user data from database
type User struct {
	ID       string `json:"id" db:"users.id"`
	Password string `json:"-" db:"users.password"`
	Email    string `json:"email" db:"users.email"`
	Username string `json:"userName" db:"users.username"`
}

// UserRole used for containg the user and the associated role
type UserRole struct {
	User
	Role
}

// Role used for creating struct for role data from database
type Role struct {
	UserID      string `json:"-" db:"userId"`
	Name        string `json:"name" db:"roles.name"`
	Permissions string `json:"-" db:"roles.permissions"`
}

// LoginDetails for the username and password coming from the request body
type LoginDetails struct {
	Username string
	Password string
}

// JwtResponse for the response when user is authenticated
type JwtResponse struct {
	Token string `json:"token"`
}

// Payload contains the payload for the jwt token
type Payload struct {
	Scope string `json:"scope"`
}

type Claims struct {
	Data Payload `json:"data"`
	jwt.RegisteredClaims
}

type GenericResponse struct {
	Message string `json:"message"`
}
