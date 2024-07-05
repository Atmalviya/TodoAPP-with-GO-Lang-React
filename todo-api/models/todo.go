package models

import (
	"time"

	"github.com/gocql/gocql"
)

type Todo struct {
	ID          gocql.UUID `json:"id"`
	UserID      gocql.UUID `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
}
