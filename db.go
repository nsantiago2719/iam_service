package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDb struct {
	db *sqlx.DB
}

// PostgresCreate makes an instance of the database
func PostgresCreate(dsn string) (*PostgresDb, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Error connection to the database:", err)
	}

	return &PostgresDb{
		db: db,
	}, nil
}

func (p *PostgresDb) getUserWithRolesByUsername(username string) ([]*UserRole, error) {
	userRoles := []*UserRole{}
	query := `
  SELECT users.id AS "users.id",
       users.email AS "users.email",
       users.username AS "users.username",
       users.password AS "users.password",
       roles.name AS "roles.name",
       roles.permissions AS "roles.permissions",
       ur.user_id AS "userId"
  FROM users AS users
  LEFT JOIN users_roles AS ur ON users.id = ur.user_id
  LEFT JOIN roles AS roles ON roles.id = ur.role_id
  WHERE users.username=$1
  `

	if err := p.db.Select(&userRoles, query, username); err != nil {
		return nil, err
	}

	return userRoles, nil
}
