package main

import "github.com/asdine/storm"

type Store struct {
	storm *storm.DB
}

func NewStore(path string) (*Store, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}
	return &Store{
		storm: db,
	}, nil
}

func (s *Store) Close() error {
	err := s.storm.Close()
	if err != nil {
		return err
	}

	return nil
}
