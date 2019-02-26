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

type Server struct {
	ID int `storm:"id,increment"`
	Name string `storm:"index"`
	Url string `storm:"unique,index"`
	Port uint16
	Timeout time.Duration
	CreatedAt time.Time
	UpdatedAt time.Time
	//ServerTypesID int
}

type Stats struct {
	ID int `storm:"id,increment"`
	MaxPlayer int64
	CurrentPlayer int64
	Time time.Time `storm:"index"`
	ServerID uint `storm:"index"`
}
