package models

import (
	"time"
)

//
//type ServerTypes struct {
//	ID int `storm:"id,increment"`
//	Kind string `storm:"unique"`
//	slug
//}

type Role int8

const AdminRole Role = 1
const ModRole Role = 1
const UserRole Role = 3

type User struct {
	ID        int    `storm:"id,increment"`
	Username  string `storm:"index,unique"`
	Hash      string `storm:"index"`
	SessionID string `storm:"index"`
	Role      Role
}

type Server struct {
	ID        int    `storm:"id,increment"`
	Name      string `storm:"index"`
	Url       string `storm:"index,unique"`
	Port      uint16
	Timeout   time.Duration
	CreatedAt time.Time
	UpdatedAt time.Time
	AddedBy   int           `storm:"index"`
	Every     time.Duration `storm:"index"`
}

type Stats struct {
	ID            int `storm:"id,increment"`
	MaxPlayer     int64
	CurrentPlayer int64
	Time          time.Time `storm:"index"`
	ServerID      uint      `storm:"index"`
}
