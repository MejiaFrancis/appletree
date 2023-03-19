package data

import (
	"time"
)

// Shool represent one row of data in our schools
type School struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Level    string    `json:"level"`
	Contact  string    `json:"contact"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	Website  string    `json:"website,omitempty"`
	Address  string    `json:"address"`
	Mode     []string  `json:"mode"`
	CreateAt time.Time `json:"-"`
	Version  int32     `json:"version"`
}
