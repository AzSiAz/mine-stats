package store

import (
	"github.com/asdine/storm"
	"mine-stats/models"
)

type Store struct {
	orm *storm.DB
}

var store *Store
var initDone = false

func NewStore(path string) (*Store, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}
	err = initStorm(db)
	if err != nil {
		return nil, err
	}

	store = &Store{orm: db}
	initDone = true

	return store, nil
}

func GetStore() *Store {
	if initDone {
		return store
	} else {
		panic("Trying to access store without init first")
	}
}

func initStorm(db *storm.DB) (err error) {
	//err = db.Init(&models.ServerTypes{})
	err = db.Init(&models.Server{})
	err = db.Init(&models.Stats{})
	err = db.Init(&models.User{})

	return
}

func (s *Store) Close() (err error) {
	err = s.orm.Close()

	return
}
